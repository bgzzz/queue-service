package queuelib

type Msg struct {
	Msg string

	// -1 if end of file
	LineNumber int
	// we can put here TTL or somthing
}

type Queue struct {
	msgs map[int]Msg

	pullingIndex int
}

func NewQueue() *Queue {

	return &Queue{
		msgs: map[int]Msg{},
	}
}

func (q *Queue) Push(m Msg) {
	q.msgs[m.LineNumber] = m
}

func (q *Queue) Pull() *Msg {

	if q.pullingIndex < 0 {
		return nil
	}

	msg, ok := q.msgs[q.pullingIndex]
	if !ok {

		result := q.msgs[-1]
		delete(q.msgs, -1)
		return &result
	}

	q.pullingIndex++

	return &msg
}
