package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/bgzzz/queue-service/pkg/queuelib"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	flagDebug        = "debug"
	flagFileParam    = "file"
	flagWriteMode    = "write"
	flagQueueName    = "queue-name"
	flagQueueService = "queue-service"
)

func main() {

	app := cli.App{
		Action: func(c *cli.Context) error {

			logger := logrus.New()
			debug := c.Bool(flagDebug)
			isWrite := c.Bool(flagDebug)

			if debug {
				logger.SetLevel(logrus.DebugLevel)
			}

			log := logrus.NewEntry(logger)

			f := c.String(flagFileParam)

			queueName := c.String(flagQueueName)

			return run(log, isWrite, f, queueName, flagQueueService)
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    flagDebug,
				Usage:   "Debug logging, ex.: (--debug)",
				EnvVars: []string{"DEBUG"},
			},
			&cli.BoolFlag{
				Name:    flagWriteMode,
				Usage:   "Specify if utility is used as subscriber, ex.: (--write)",
				Value:   false,
				EnvVars: []string{"WRITE_MODE"},
			},
			&cli.StringFlag{
				Name:    flagFileParam,
				Usage:   "Specify file to read or name of the file to write, ex (--file input-output.txt)",
				EnvVars: []string{"FILE_PATH"},
			},
			&cli.StringFlag{
				Name:    flagQueueName,
				Usage:   "Specify queue name to read from/write to  (--queue-name file-exchange)",
				EnvVars: []string{"QUEUE_NAME"},
			},
			&cli.StringFlag{
				Name:    flagQueueService,
				Usage:   "Specify queue service name to read from/write to  (--queue-service http://localhost:8080)",
				EnvVars: []string{"QUEUE_NAME"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func run(log *logrus.Entry, isWrite bool, f, qName, qServiceName string) error {
	// validate queue connection
	resp, err := http.Get(path.Join(qServiceName, "queues",
		uuid.NewString()))
	if err != nil {
		return errors.Wrap(err, "unable to check queue connection")
	}

	if resp.StatusCode != http.StatusBadRequest {
		return errors.New("queue is not working")
	}

	if isWrite {

		file, err := os.Create(f)
		if err != nil {
			return err
		}
		defer file.Close()

		w := bufio.NewWriter(file)
		for {
			resp, err := http.Get(path.Join(qServiceName, "queues",
				qName))
			if err != nil {
				return errors.Wrap(err, "pull from queue")
			}

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return errors.Wrap(err, "unable to read body")
			}

			var msg queuelib.Msg
			if err := json.Unmarshal(b, &msg); err != nil {
				return errors.Wrap(err, "unable to marshal response")
			}

			log.Info("read line %v from queue %v", msg.Msg, qName)
			fmt.Fprintln(w, msg.Msg)

			if msg.LineNumber == -1 {
				break
			}
		}
		w.Flush()
		return nil
	}

	// validating file
	file, err := os.Open(f)
	if err != nil {
		return errors.Wrap(err, "unable to open file for reading")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	indx := 0
	for scanner.Scan() {
		line := scanner.Text()
		log.Infof("Reading line %v", line)

		if err := sendMsg(line, qServiceName, qName, indx); err != nil {
			return err
		}
		indx++
	}

	if err := sendMsg("", qServiceName, qName, -1); err != nil {
		return err
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func sendMsg(line, qServiceName, qName string, indx int) error {
	msg := queuelib.Msg{
		Msg:        line,
		LineNumber: indx,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "unable to marshal request")
	}

	resp, err := http.Post(path.Join(qServiceName, qName),
		"application/json", bytes.NewBuffer(b))
	if err != nil {
		return errors.Wrap(err, "unable to push message to the queue")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(err, "something wrong with the queue ")
	}

	return nil
}
