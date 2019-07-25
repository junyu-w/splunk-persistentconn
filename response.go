package persistentconn

// Response represents the response sent back to the client
type Response struct {
	statusCode int
	body       string
	slotIndex  int
}

type responseQueueSlot struct {
	resp Response
}

type responseQueue []responseQueueSlot

func newResponseQueue() responseQueue {
	return make([]responseQueueSlot, 0)
}
