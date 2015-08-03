package gotcp

import (
	"net"
	"sync"
	"time"
)

type Config struct {
	PacketSendChanLimit    uint32 // the limit of packet send channel
	PacketReceiveChanLimit uint32 // the limit of packet receive channel
	WorkerNum              uint32
}

type Server struct {
	config    *Config         // server configuration
	callback  ConnCallback    // message callbacks in connection
	protocol  Protocol        // customize packet protocol
	exitChan  chan struct{}   // notify all goroutines to shutdown
	waitGroup *sync.WaitGroup // wait for all goroutines
}

// NewServer creates a server
func NewServer(config *Config, callback ConnCallback, protocol Protocol) *Server {
	return &Server{
		config:    config,
		callback:  callback,
		protocol:  protocol,
		exitChan:  make(chan struct{}),
		waitGroup: &sync.WaitGroup{},
	}
}

// Start starts service
func (s *Server) Start(listener *net.TCPListener, acceptTimeout time.Duration) {
	s.waitGroup.Add(1)
	defer func() {
		listener.Close()
		s.waitGroup.Done()
	}()

	for {
		select {
		case <-s.exitChan:
			return

		default:
		}

		listener.SetDeadline(time.Now().Add(acceptTimeout))

		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}

		if s.config.WorkerNum == 0 {
			go newConn(conn, s).Do()
		} else {
			go newConn(conn, s).DoPool(s.config.WorkerNum)
		}

	}
}
func (s *Server) StartConnector(reConnect time.Duration) {
	s.waitGroup.Add(1)
	defer func() {
		s.waitGroup.Done()
	}()

	for {
		conn, err := s.callback.OnDial()
		if err != nil {
			select {
			case <-s.exitChan:
				return
			case <-time.After(time.Second * reConnect):
				continue
			}
		}

		c := newConn(conn, s)
		go c.Do()
		if s.config.WorkerNum == 0 {
			go c.Do()
		} else {
			go c.DoPool(s.config.WorkerNum)
		}

		select {
		case <-s.exitChan:
			return
		case <-c.closeChan:
			time.Sleep(time.Second * reConnect)
		}
	}

}

// Stop stops service
func (s *Server) Stop() {
	close(s.exitChan)
	s.waitGroup.Wait()
}
