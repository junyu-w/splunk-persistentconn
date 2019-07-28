package persistentconn

import (
	"encoding/json"
	"fmt"
)

// Response represents the response sent back to the client
type Response struct {
	StatusCode int    `json:"status"`
	Body       string `json:"payload"`
}

// getRawData transforms the response to a payload that splunkd can decode
func (resp Response) getRawData() string {
	respData, err := json.Marshal(&resp)
	if err != nil {
		respData = []byte("Failed to serialize response data")
	}
	rawData := fmt.Sprintf("%d\n%s", len(respData), respData)
	return rawData
}
