package queuelib

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddRemoveFromQueue(t *testing.T) {
	tests := []struct {
		initMsgs     []Msg
		addMsgs      []Msg
		rmMsgs       int
		expectedMsgs []Msg
	}{
		{
			initMsgs: []Msg{},
			addMsgs: []Msg{
				{
					Msg:        "add",
					LineNumber: 1,
				},
			},
			rmMsgs: 0,
			expectedMsgs: []Msg{
				{
					Msg:        "add",
					LineNumber: 1,
				},
			},
		},

		{
			initMsgs: []Msg{
				{
					Msg:        "add",
					LineNumber: 1,
				},
			},
			addMsgs:      []Msg{},
			rmMsgs:       1,
			expectedMsgs: []Msg{},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test_%d", i),
			func(t *testing.T) {

				q := NewQueue()

				for _, m := range test.initMsgs {
					q.Push(m)
				}

				for _, m := range test.addMsgs {
					q.Push(m)
				}

				for i := 0; i < test.rmMsgs; i++ {
					q.Pull()
				}

				assert.Equal(t, len(test.expectedMsgs), len(q.msgs),
					"should be the same length")

			})
	}
}
