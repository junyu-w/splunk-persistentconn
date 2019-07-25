package persistentconn

// Response represents the response sent back to the client
type Response struct {
	statusCode int
	body       string
}

type responseQueueSlot struct {
	resp *Response
}

type responseQueue []responseQueueSlot

func newResponseQueue() responseQueue {
	return make([]responseQueueSlot, 0)
}

func (rq responseQueue) allocateNewSlot() responseQueueSlot {
	slot := responseQueueSlot{}
	rq = append(rq, slot)
	return slot
}
