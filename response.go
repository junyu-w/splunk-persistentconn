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
// splunkd response protocol: 0\n<len_response_bytes>\n<response>
// NOTE: leading 0 in protocol is necessary to let splunkd know it's a response
func (resp Response) getRawData() string {
	respData, err := json.Marshal(&resp)
	if err != nil {
		respData = []byte("Failed to serialize response data")
	}
	rawData := fmt.Sprintf("0\n%d\n%s", len(respData), respData)
	return rawData
}
