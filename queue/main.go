package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/bgzzz/queue-service/pkg/queuelib"
)

func handleQueues(w http.ResponseWriter, req *http.Request,
	queueSrv *QueueServer, log *logrus.Entry) {

	qName, err := extractQueueName(req.URL.Path)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Method == http.MethodPost {

		msg := queuelib.Msg{}
		if err := json.NewDecoder(req.Body).Decode(&msg); err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		queueSrv.AddToQueue(qName, &msg)
		return
	}

	if req.Method != http.MethodGet {
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	msg, err := queueSrv.RemoveFromQueue(qName)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(msg)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(js); err != nil {
		log.Error(err)
	}
}

func extractQueueName(url string) (string, error) {
	tmp := strings.Split(url, "/")

	if len(tmp) != 3 {
		return "", errors.New("unable to find queue name")
	}

	if tmp[1] != "queues" {
		return "", errors.New("unsupported resource")
	}

	return tmp[2], nil
}

func main() {

	queueSrv := NewQueueSrv()

	logger := logrus.New()

	log := logrus.NewEntry(logger)

	if err := http.ListenAndServe(":8090", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handleQueues(w, req, queueSrv, log)
	})); err != nil {
		log.Error(err)
	}
}
