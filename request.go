package persistentconn

import "encoding/json"

// Request represents the request sent from client
type Request struct {
	OutputMode         string `json:"output_mode"`
	OutputModeExplicit bool   `json:"output_mode_explicit"`
	Server             struct {
		RestURI    string `json:"rest_uri"`
		Hostname   string `json:"hostname"`
		Servername string `json:"servername"`
		GUID       string `json:"guid"`
	} `json:"server"`
	Restmap struct {
		Name string `json:"name"`
		Conf struct {
			Handler         string `json:"handler"`
			Match           string `json:"match"`
			OutputModes     string `json:"output_modes"`
			PassHTTPHeaders string `json:"passHttpHeaders"`
			PassPayload     string `json:"passPayload"`
			Script          string `json:"script"`
			Scripttype      string `json:"scripttype"`
		} `json:"conf"`
	} `json:"restmap"`
	PathInfo   string     `json:"path_info"`
	Query      [][]string `json:"query"`
	Connection struct {
		SrcIP         string `json:"src_ip"`
		Ssl           bool   `json:"ssl"`
		ListeningPort int    `json:"listening_port"`
	} `json:"connection"`
	Session struct {
		User      string `json:"user"`
		Authtoken string `json:"authtoken"`
	} `json:"session"`
	RestPath string `json:"rest_path"`
	Method   string `json:"method"`
	Ns       struct {
		App  string `json:"app"`
		User string `json:"user"`
	} `json:"ns"`
	Headers [][]string `json:"headers"`
	Form    [][]string `json:"form,omitempty"`
	Payload string     `json:"payload,omitempty"`
}

// parseRequests creates a Request object by parsing information from a request packet.
func parseRequest(p *RequestPacket) (Request, error) {
	block := p.block
	var request Request
	if err := json.Unmarshal([]byte(block), &request); err != nil {
		return Request{}, err
	}
	return request, nil
}
