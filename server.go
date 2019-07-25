package persistentconn

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Server represents the persistentconn server that handles request
// and writes response back to the client
type Server struct {
	requestChan   chan Request
	responseChan  chan responseQueueSlot
	responseQueue responseQueue
	registry      *handlerRegistry
}

// NewServer creates a persistentconn server
func NewServer() *Server {
	return &Server{
		requestChan:   make(chan Request),
		responseChan:  make(chan responseQueueSlot),
		responseQueue: newResponseQueue(),
		// TODO: add registry initialization
		registry: nil,
	}
}

// Run starts a persistentconn server and starts handling request sent from
// client (with splunkd as the middle layer)
func (s *Server) Run() {
	go s.readInputPacket()
	go s.processResponse()
	for req := range s.requestChan {
		slot := s.responseQueue.allocateNewSlot()
		handler := s.registry.getHandler(req.PathInfo)
		// handle request in a goroutine
		go func(req Request, respSlot responseQueueSlot) {
			resp, err := handler(req)
			if err != nil {
				resp = Response{
					statusCode: http.StatusInternalServerError,
					body:       err.Error(),
				}
			}
			slot.resp = &resp
			s.responseChan <- respSlot
		}(req, slot)
	}
}

// readInputPacket starts a separate goroutine that reads request sent from client
// and is the entrypoint of a server process
func (s *Server) readInputPacket() {
	for {
		inPacket, err := ReadPacket(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		req, err := parseRequest(inPacket)
		if err != nil {
			log.Fatal(err)
		}
		s.requestChan <- req
	}
}

// processResponse proccesses response from handler and sent the response back to the client
func (s *Server) processResponse() {
	for slot := range s.responseChan {
		// TODO: implement response queue checking and writing response to stdout
		fmt.Printf("Got response - status: %d - body: %s\n", slot.resp.statusCode, slot.resp.body)
	}
}
