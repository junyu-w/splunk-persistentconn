package persistentconn

import (
	"fmt"
)

// Response represents the response sent back to the client
type Response struct {
	StatusCode int
	Body       string
}

// getRawData transforms the response to a payload that splunkd can decode
func (resp Response) getRawData() string {
	rawData := fmt.Sprintf("%d\n%s", len(resp.Body), resp.Body)
	return rawData
}
