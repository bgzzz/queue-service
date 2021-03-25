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

	if req.Method == http.MethodPost {
		qName, err := extractQueueName(req.URL.Path)
		if err != nil {
			log.Error(err)
			http.Error(w, "queue is unspecified in url", http.StatusBadRequest)
			return
		}

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

	qName, err := extractQueueName(req.URL.Path)
	if err != nil {
		log.Error(err)
		http.Error(w, "queue is unspecified in url", http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func extractQueueName(url string) (string, error) {
	tmp := strings.Split(url, "/")

	if len(tmp) != 2 {
		return "", errors.New("unable to find queue name")
	}

	return tmp[1], nil
}

func (qs *QueueServer) RemoveFromQueue(qName string) (*queuelib.Msg, error) {
	qs.mtx.Lock()
	defer qs.mtx.RUnlock()

	q, ok := qs.Queues[qName]
	if !ok {
		return nil, errors.New("queue does not exist")
	}

	msg := q.Pull()
	if msg == nil {
		return nil, errors.New("queue is empty")
	}

	return msg, nil

}

func main() {

	queueSrv := NewQueueSrv()

	logger := logrus.New()

	log := logrus.NewEntry(logger)

	http.HandleFunc("/queues", func(w http.ResponseWriter, req *http.Request) {
		handleQueues(w, req, queueSrv, log)
	})

	http.ListenAndServe(":8090", nil)
}
