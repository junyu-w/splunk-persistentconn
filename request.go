package persistentconn

import (
	"encoding/json"
	"log"
)

// splunkdRequest represents the request sent from splunkd
type splunkdRequest struct {
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

// Request contains information of an incoming request
type Request struct {
	// isInit indicates if this request corresponds to a packet with command byte set to INIT
	// NOTE: even though a packet's command byte could be both INIT and BLOCK, we will separate them out
	// into two different Request objects to avoid confusion. Request object with isInit: true should not
	// have other attributes set (i.e shouldn't be a data request and an init request at the same time).
	isInit     bool
	OutputMode string            `json:"output_mode"`
	Headers    map[string]string `json:"headers"`
	Method     string            `json:"method"`
	Namespace  struct {
		App  string `json:"app"`
		User string `json:"user"`
	} `json:"namespace"`
	Session struct {
		User      string `json:"user"`
		Authtoken string `json:"authtoken"`
	} `json:"session"`
	Query   map[string]string `json:"query"`
	Form    map[string]string `json:"form,omitempty"`
	Payload string            `json:"payload,omitempty"`
	Path    string            `json:"path"`
	Params  map[string]string `json:"params,omitempty"`
}

// parseRequests creates a Request object by parsing information from a request packet.
func (s *Server) parseRequest(p *RequestPacket) {
	if p.isFirst() {
		request := Request{isInit: true}
		s.requestChan <- request
	}
	if p.hasBlock() {
		block := p.block
		var splunkdReq splunkdRequest
		if err := json.Unmarshal([]byte(block), &splunkdReq); err != nil {
			log.Fatal(err)
		}
		request := Request{
			OutputMode: splunkdReq.OutputMode,
			Headers:    tupleListToMap(splunkdReq.Headers),
			Method:     splunkdReq.Method,
			Namespace:  splunkdReq.Ns,
			Session:    splunkdReq.Session,
			Query:      tupleListToMap(splunkdReq.Query),
			Form:       tupleListToMap(splunkdReq.Form),
			Payload:    splunkdReq.Payload,
			Path:       splunkdReq.PathInfo,
			Params:     make(map[string]string),
		}
		s.requestChan <- request
	}
}
