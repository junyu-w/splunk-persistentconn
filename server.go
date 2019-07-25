package persistentconn

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Server represents the persistentconn server that handles request
// and writes response back to the client
type Server struct {
	requestChan   chan Request
	responseChan  chan Response
	responseQueue responseQueue
	registry      *handlerRegistry
}

// NewServer creates a persistentconn server
func NewServer() *Server {
	return &Server{
		requestChan:   make(chan Request),
		responseChan:  make(chan Response),
		responseQueue: newResponseQueue(),
		// TODO: add registry initialization
		registry: &handlerRegistry{},
	}
}

// Run starts a persistentconn server and starts handling request sent from
// client (with splunkd as the middle layer)
func (s *Server) Run() {
	go s.handleRequest()
	go s.processResponse()
	s.startProcessingInputPackets()
}

// readInputPacket starts a separate goroutine that reads request sent from client
// and is the entrypoint of a server process
func (s *Server) startProcessingInputPackets() {
	for {
		inPacket, err := ReadPacket(os.Stdin)
		if err != nil {
			if err == io.EOF {
				continue
			}
			log.Fatal(err)
		}
		req, err := parseRequest(inPacket)
		if err != nil {
			log.Fatal(err)
		}
		s.requestChan <- req
	}
}

func (s *Server) handleRequest() {
	for req := range s.requestChan {
		s.responseQueue = append(s.responseQueue, responseQueueSlot{})
		handler := s.registry.getHandler(req.PathInfo)
		// handle request in a goroutine
		go func(req Request, respIndex int) {
			resp, err := handler(req)
			if err != nil {
				resp = Response{
					statusCode: http.StatusInternalServerError,
					body:       err.Error(),
				}
			}
			resp.slotIndex = respIndex
			s.responseChan <- resp
		}(req, len(s.responseQueue)-1)
	}
}

// processResponse proccesses response from handler and sent the response back to the client
func (s *Server) processResponse() {
	for resp := range s.responseChan {
		// TODO: implement response queue checking and writing response to stdout
		fmt.Printf("Got response - status: %d - body: %s - index: %d\n", resp.statusCode, resp.body, resp.slotIndex)
	}
}
