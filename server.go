package persistentconn

import (
	"bufio"
	"container/list"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

// Server represents the persistentconn server that handles request
// and writes response back to the client
type Server struct {
	requestChan       chan Request
	responseChan      chan Response
	responseQueue     *list.List
	registry          *handlerRegistry
	resposneQueueLock *sync.Mutex
}

// NewServer creates a persistentconn server
func NewServer() *Server {
	return &Server{
		requestChan:       make(chan Request),
		responseChan:      make(chan Response),
		responseQueue:     list.New(),
		registry:          &handlerRegistry{},
		resposneQueueLock: new(sync.Mutex),
	}
}

// Handle registers a handler function for a given path (or path pattern).
// A path pattern is in the format of "<component>/:<param_1>/<component>/..." and a path component
// starting with ":" indicates it's a parameter which will be inferred from the actual path in the request
// E.g. if the registered path pattern is "entity/:name/data" and the
// path in the request is "entity/hello/data", then the key-value pair {"name": "hello"} will be stored
// in the request's params which can be later referenced inside of the handler.
func (s *Server) Handle(path string, handler Handler, allowedMethods ...string) {
	s.registry.register(path, handler, allowedMethods)
}

// Run starts a persistentconn server and starts handling request sent from
// client (with splunkd as the middle layer)
// FIXME: splunkd always starts new server process upon request
func (s *Server) Run() {
	go s.handleRequest()
	go s.processResponse()
	s.startProcessingInputPackets(os.Stdin)
}

// startProcessingInputPackets starts a separate goroutine that reads request sent from client
// and is the entrypoint of a server process
func (s *Server) startProcessingInputPackets(input io.Reader) {
	for {
		inPacket, err := ReadPacket(input)
		if err != nil {
			if err == io.EOF {
				continue
			}
			log.Fatal(err)
		}
		if inPacket.hasBlock() {
			req, err := parseRequest(inPacket)
			if err != nil {
				log.Fatal(err)
			}
			s.requestChan <- req
		}
	}
}

// handleRequest takes request that comes in and find the corresponding handler
func (s *Server) handleRequest() {
	for req := range s.requestChan {
		s.resposneQueueLock.Lock()
		elem := s.responseQueue.PushBack(struct{}{})
		s.resposneQueueLock.Unlock()

		// handle request in a goroutine
		handler := s.registry.getHandler(req)
		go func(req Request, slot *list.Element) {
			resp, err := handler(req)
			if err != nil {
				resp = Response{
					StatusCode: http.StatusInternalServerError,
					Body:       err.Error(),
				}
			}
			// TODO: replace all print statements with proper logging
			// fmt.Printf("Finished handling - response - status: %d - body: %s\n", resp.StatusCode, resp.Body)
			slot.Value = resp
			s.responseChan <- resp
		}(req, elem)
	}
}

// processResponse proccesses response from handler and sent the response back to the client
func (s *Server) processResponse() {
	for range s.responseChan {
		flushedCount, err := s.flushResponses(os.Stdout)
		if err != nil {
			// fmt.Println("Failed to flush response - Error:", err)
			continue
		}
		if flushedCount != 0 {
			// fmt.Printf("Flushed %d responses\n", flushedCount)
		}
	}
}

// flushResponses go through responses in the response queue of the server, and it flushes consecutive
// responses starting from the front of the queue in batch to ensure that responses are synchronized in the same
// order as the corresonding requests.
func (s *Server) flushResponses(output io.Writer) (int, error) {
	s.resposneQueueLock.Lock()
	defer s.resposneQueueLock.Unlock()
	// prepare response data to flush to stdout
	elem := s.responseQueue.Front()
	flushedElList := make([]*list.Element, 0)

	writer := bufio.NewWriter(output)
	for {
		if elem == nil {
			break
		}
		resp, ok := elem.Value.(Response)
		if !ok {
			break
		}
		data := resp.getRawData()
		_, err := writer.WriteString(data)
		if err != nil {
			return 0, err
		}
		writer.WriteString("\n")
		flushedEl := elem
		flushedElList = append(flushedElList, flushedEl)
		elem = elem.Next()
	}
	err := writer.Flush()
	if err != nil {
		return 0, err
	}
	// clean up flushed element from the queue
	for _, flushedEl := range flushedElList {
		s.responseQueue.Remove(flushedEl)
	}
	return len(flushedElList), nil
}
