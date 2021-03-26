package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAddRemoveFromQueue(t *testing.T) {

	queueSrv := NewQueueSrv()

	logger := logrus.New()

	log := logrus.NewEntry(logger)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handleQueues(w, req, queueSrv, log)
	}))
	defer ts.Close()

	tests := []struct {
		reqURL             string
		payload            string
		httpMethod         string
		expectedHttpStatus int
		expectedPayload    string
	}{
		{
			reqURL:             "queues/not-existing-queue",
			payload:            "",
			httpMethod:         "GET",
			expectedHttpStatus: http.StatusBadRequest,
			expectedPayload:    "",
		},
		{
			reqURL:             "queues/rw",
			payload:            `{ "msg": "hello", "lineNumber": 1}`,
			httpMethod:         "POST",
			expectedHttpStatus: http.StatusOK,
			expectedPayload:    "",
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test_%d", i),
			func(t *testing.T) {
				var res *http.Response
				var err error
				if test.httpMethod == "GET" {
					res, err = http.Get(fmt.Sprintf("%s/%s",
						ts.URL, test.reqURL))
				} else {
					//TODO
					return
				}
				if err != nil {
					t.Errorf("error %v", err)
				}

				assert.Equal(t, test.expectedHttpStatus, res.StatusCode)
				if test.expectedPayload != "" {
					//TODO
					return
				}
			})
	}
}
