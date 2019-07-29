package persistentconn

import (
	"encoding/json"
	"fmt"
)

// Response represents the response sent back to the client
type Response struct {
	isInit     bool
	StatusCode int    `json:"status"`
	Body       string `json:"payload"`
}

// getRawData transforms the response to a payload that splunkd can decode
// splunkd protocol for response to init packet: "0\n" (empty byte with length 0) to indicate success
// splunkd protocl for response to data packet: <len_response_bytes>\n<response>
func (resp Response) getRawData() string {
	if resp.isInit {
		return "0\n"
	}
	respData, err := json.Marshal(&resp)
	if err != nil {
		respData = []byte("Failed to serialize response data")
	}
	rawData := fmt.Sprintf("%d\n%s", len(respData), respData)
	return rawData
}
