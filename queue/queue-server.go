package main

import (
	"errors"
	"sync"

	"github.com/bgzzz/queue-service/pkg/queuelib"
)

type QueueServer struct {
	mtx    sync.RWMutex
	Queues map[string]*queuelib.Queue
}

func NewQueueSrv() *QueueServer {
	return &QueueServer{
		Queues: map[string]*queuelib.Queue{},
	}
}

func (qs *QueueServer) AddToQueue(qName string, msg *queuelib.Msg) {
	qs.mtx.Lock()
	defer qs.mtx.Unlock()

	q, ok := qs.Queues[qName]
	if !ok {
		qs.Queues[qName] = queuelib.NewQueue()
		qs.Queues[qName].Push(*msg)
		return
	}

	q.Push(*msg)
}

func (qs *QueueServer) RemoveFromQueue(qName string) (*queuelib.Msg, error) {
	qs.mtx.RLock()
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
