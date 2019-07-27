package persistentconn

import (
	"bufio"
	"fmt"
	"os"
)

// Response represents the response sent back to the client
type Response struct {
	StatusCode    int
	Body          string
	isPlaceholder bool
}

func (resp Response) getRawData() string {
	rawData := fmt.Sprintf("%d\n%s", len(resp.Body), resp.Body)
	return rawData
}

type responseQueue []Response

func newResponseQueue() responseQueue {
	return make([]Response, 0)
}

func (rq responseQueue) flushResponses() (int, error) {
	flushedCount := 0
	writer := bufio.NewWriter(os.Stdout)
	for _, resp := range rq {
		if resp.isPlaceholder {
			break
		}
		data := resp.getRawData()
		_, err := writer.WriteString(data)
		if err != nil {
			return 0, err
		}
		writer.WriteString("\n")
		flushedCount++
	}

	err := writer.Flush()
	if err != nil {
		return 0, err
	}
	return flushedCount, nil
}
