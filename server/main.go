package server

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/krysopath/klystron/structs"
	"github.com/krysopath/klystron/utils"
)

type Server struct {
	//Logger     *log.Logger
	SocketAddr       string
	SocketBufferSize int
}

func (s *Server) consumeSocket(c net.Conn) []byte {
	received := make([]byte, 0)
	log.Printf("Receiving with %d buffer...", s.SocketBufferSize)
	for {
		buf := make([]byte, s.SocketBufferSize)
		count, err := c.Read(buf)
		received = append(received, buf[:count]...)
		if err != nil {
			if err != io.EOF {
				log.Fatalf("Error on read: %s", err)
			}
			break
		}
	}
	return received
}

func handleSources(job *structs.Job) {
	if _, err := os.Stat(job.Directory); os.IsNotExist(err) {
		os.Mkdir(job.Directory, 0700)
	}
	for _, source := range job.Sources {
		//data := []byte(source.Data)
		dataDecoded, _ := base64.StdEncoding.DecodeString(source.Data)
		//data := utils.GUnzipData([]byte(source.Data))
		var filePath string = fmt.Sprintf("%s/%s", structs.JobSpoolDir, source.Hash)
		if !utils.FileExists(filePath) {
			ioutil.WriteFile(filePath, dataDecoded, 0600)
			log.Printf("Wrote source: %s@%s", source.Path, source.Hash)
		}
	}
}

func handleJob(job *structs.Job) {
	log.Printf("Got Job: %s", job.Name)
	log.Printf("Using Spool Directory: %s", structs.JobSpoolDir)
	log.Printf("Outputs: %+v", job.Outputs)
	handleSources(job)

}

func handleMessage(messageBytes *[]byte) {
	//message := utils.JSONUnmarshal(messageBytes).(map[string]interface{})
	message := utils.JSONUnmarshal(messageBytes).(structs.Message)
	fmt.Printf("%+v", message)
	handleJob(&message.Content)
	//pdf.CreateCsv("examples/csv/addresses.csv")

}

func (s *Server) handleConn(c net.Conn) {
	log.Printf("Got Conn: %+v", c)
	defer c.Close()
	received := s.consumeSocket(c)
	handleMessage(&received)
	c.Close()
}

func (s *Server) Listen() {
	listener, err := net.Listen("unix", s.SocketAddr)
	if err != nil {
		log.Fatalf("Unable to listen on socket file %s: %s",
			s.SocketAddr, err)
	}
	defer listener.Close()
	log.Printf("klystron bound to %s", s.SocketAddr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error on accept: %s", err)
		}
		go s.handleConn(conn)
	}
}

func NewServer(socketAddr string) Server {
	os.Remove(socketAddr)
	return Server{
		SocketAddr:       socketAddr,
		SocketBufferSize: structs.SockBufferSize64,
	}
}
