package main

import (
	"fmt"
	"log"
	"net"
)

// Server struct
type Server struct {
	listenAddress string
	ln            net.Listener
	quitch        chan struct{}
}

// Newserver function
func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddress: listenAddr,
		quitch:        make(chan struct{}),
	}
}

// Make accept loop
func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Accept error: ", err)
			continue
		}

		fmt.Println("New connection to the server: ", conn.RemoteAddr())
		go s.readLoop(conn)
	}
}

// Make read loop
func (s *Server) readLoop(con net.Conn) {

	defer con.Close()

	// Make bytes
	buf := make([]byte, 2048)
	for {
		n, err := con.Read(buf)
		if err != nil {
			fmt.Println("Read error: ", err)
			continue
		}

		msg := buf[:n]
		fmt.Println(string(msg))
	}
}

// start server
func (s *Server) Start() error {

	ln, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		return err
	}

	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()

	<-s.quitch

	return nil

}

func main() {
	server := NewServer(":3001")
	log.Fatal(server.Start())
}
