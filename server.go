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

// startProcessingInputPackets starts a separate goroutine that reads request sent from client
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
		// TODO: we should only process packet that has a input block. opcode=0x02
		req, err := parseRequest(inPacket)
		if err != nil {
			log.Fatal(err)
		}
		s.requestChan <- req
	}
}

func (s *Server) handleRequest() {
	for req := range s.requestChan {
		s.responseQueue = append(s.responseQueue, Response{isPlaceholder: true})
		handler := s.registry.getHandler(req.PathInfo)
		// handle request in a goroutine
		go func(req Request, respIndex int) {
			resp, err := handler(req)
			if err != nil {
				resp = Response{
					StatusCode: http.StatusInternalServerError,
					Body:       err.Error(),
				}
			}
			fmt.Printf("Got response - status: %d - body: %s - index: %d\n", resp.StatusCode, resp.Body, respIndex)
			// FIXME: race condition where flushing has shrinked the queue so resulting in index out of range :(
			s.responseQueue[respIndex] = resp
			s.responseChan <- resp
		}(req, len(s.responseQueue)-1)
	}
}

// processResponse proccesses response from handler and sent the response back to the client
func (s *Server) processResponse() {
	for range s.responseChan {
		flushedCount, err := s.responseQueue.flushResponses()
		if err != nil {
			fmt.Println("Failed to flush response - Error:", err)
			continue
		}
		fmt.Printf("Flushed %d responses\n", flushedCount)
		s.responseQueue = s.responseQueue[flushedCount:]
	}
}
