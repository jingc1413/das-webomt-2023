package udp

import (
	"gomt/core/net/pool"
	"net"
	"strings"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type MessageHandler interface {
	HandleMessage(addr string, payload []byte)
}

type Message struct {
	Addr    string
	Payload []byte
}

type Server struct {
	id   string
	name string
	log  *log.Entry

	host      string
	pool      pool.WorkerPool
	handler   MessageHandler
	limit     chan bool
	conn      *net.UDPConn
	sendQueue chan Message

	wg       sync.WaitGroup
	quit     chan bool
	shutdown bool
}

func NewServer(id string, host string, handler MessageHandler, connLimit int, maxWorkers int) (*Server, error) {
	s := &Server{
		name:      "udp",
		id:        id,
		host:      host,
		quit:      make(chan bool),
		handler:   handler,
		limit:     make(chan bool, connLimit),
		sendQueue: make(chan Message),
	}
	s.log = log.WithFields(log.Fields{"udp-server": s.id})
	s.pool = pool.NewWorkerPool("udp", maxWorkers, maxWorkers*100)
	return s, nil
}

func (s Server) String() string {
	return s.name + ":" + s.id
}

func (s *Server) Run() {
	s.log.Trace("udp server is running")
	defer s.log.Trace("udp server is stopped")
	defer s.wg.Wait()

	s.log.Infof("listen on %s", s.host)
	addr, err := net.ResolveUDPAddr("udp", s.host)
	if err != nil {
		s.log.Fatal(errors.Wrap(err, "resolve address error"))
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		s.log.Fatal(errors.Wrap(err, "listen to addr error"))
		return
	}
	defer conn.Close()
	s.conn = conn

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.pool.Run()
	}()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			if s.shutdown {
				break
			}
			select {
			case s.limit <- true:
				s.wg.Add(1)
				go func() {
					defer s.wg.Done()
					s.handleConnection(conn)
				}()
			case m := <-s.sendQueue:
				go s.sendMessage(conn, m.Addr, m.Payload)
			}
		}
	}()

	<-s.quit

	s.pool.Stop()
}

func (s *Server) Stop() {
	s.log.Trace("stop udp server")
	s.shutdown = true
	s.quit <- true
}

func (s *Server) handleConnection(conn *net.UDPConn) {
	buffer := make([]byte, 2048)
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		if !s.shutdown && !strings.Contains(err.Error(), "use of closed network connection") {
			s.log.Error(errors.Wrap(err, "read from udp error"))
		}
	} else {
		job := pool.Job{
			Payload: Message{
				Addr:    addr.String(),
				Payload: buffer[:n],
			},
			HandleFunc: s.handleMessage,
		}
		s.pool.PushJob(job)
	}
	<-s.limit
}

func (s *Server) sendMessage(conn *net.UDPConn, addr string, payload []byte) {
	dst, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		s.log.Error(errors.Wrap(err, "resolve address error"))
		return
	}
	// s.log.Tracef("send: %v, %v", hex.EncodeToString(payload), addr)
	if _, err := conn.WriteTo(payload, dst); err != nil {
		s.log.Error(errors.Wrap(err, "send message error"))
		return
	}
}

func (s *Server) handleMessage(payload interface{}) {
	msg, ok := payload.(Message)
	if !ok {
		s.log.Error(errors.New("unknow payload type, handle message error"))
		return
	}
	// s.log.Tracef("recv: %v, %v", hex.EncodeToString(msg.Payload), msg.Addr)
	s.handler.HandleMessage(msg.Addr, msg.Payload)
	return
}

func (s *Server) SendMessage(addr string, payload []byte) error {
	s.sendQueue <- Message{
		Addr:    addr,
		Payload: payload,
	}
	return nil
}
