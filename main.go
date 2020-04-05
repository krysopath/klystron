package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const sockAddrDefault = "/tmp/klystron.sock"
const sockBufferSize64 = 212992
const sockBufferSize32 = 163840

type fontConfig struct {
	Name    string `json:"name"`
	Style   string `json:"style"`
	Size    int    `json:"size"`
	FontDir string `json:"fontDir"`
}

type pdfFile struct {
	Orientation string     `json:"orientation"`
	Unit        string     `json:"unit"`
	Format      string     `json:"format"`
	Font        fontConfig `json:"font"`
}

type Job struct {
	Directory string    `json:"directory"`
	Outputs   []pdfFile `json:"outputs"`
}

func validateJob(job *Job) bool {
	valid := true
	dir := job.Directory
	if string(dir) != dir {
		log.Fatal("job has weird base directory")
		valid = false
	}
	if len(job.Directory) < 1 {
		log.Println("job has weird base directory")
		valid = false
	}
	if len(job.Outputs) < 1 {
		log.Println("job has no specified outputs")
		valid = false
	}
	return valid

}

func dumpJobOpts(job *Job) {
	bytes, err := json.Marshal(job)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

func HandleConn(c net.Conn) {
	received := make([]byte, 0)
	log.Printf("Receiving...")
	for {
		buf := make([]byte, 512)
		count, err := c.Read(buf)
		received = append(received, buf[:count]...)
		if err != nil {
			var job Job
			json.Unmarshal(received, &job)
			if validateJob(&job) {
				log.Printf("Got Job: %+v", job)
				c.Write([]byte(`OK`))
			}
			c.Close()
			if err != io.EOF {
				log.Fatalf("Error on read: %s", err)
			}
			break
		}
	}
}

type Server struct {
	//Logger     *log.Logger
	SocketAddr       string
	SocketBufferSize int
}
type Sender struct {
	//Logger     *log.Logger
	SocketAddr string
}

func (s *Sender) SendMessage(message []byte) {
	c, err := net.Dial("unix", s.SocketAddr)
	if err != nil {
		log.Printf("Failed to dial: %s", err)
	}
	defer c.Close()
	count, err := c.Write(message)
	if err != nil {
		log.Printf("Write error: %s", err)
	}
	log.Printf("Wrote %d bytes", count)
}

func (s *Server) Listen() {
	listener, err := net.Listen("unix", s.SocketAddr)
	if err != nil {
		log.Fatalf("Unable to listen on socket file %s: %s", s.SocketAddr, err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error on accept: %s", err)
		}
		log.Printf("Got Conn: %+v", conn)
		go HandleConn(conn)
	}
}

func NewServer() Server {
	server := Server{}
	server.SocketAddr = sockAddrDefault
	server.SocketBufferSize = sockBufferSize64
	os.Remove(server.SocketAddr)
	return server
}

//func getStdin() *bufio.Reader {}

func (s *Sender) inputChannel() *bufio.Reader {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		fmt.Println("no pipe :(")
	} else {
		fmt.Println("hi pipe!")
	}

	if fi.Mode()&os.ModeCharDevice != 0 || fi.Size() <= 0 {
		//return sender, errors.New("Cant work without pipe")
	}
	reader := bufio.NewReader(os.Stdin)
	return reader

}

func NewSender(socketAddr string) Sender {
	sender := Sender{}
	sender.SocketAddr = socketAddr
	return sender
}

func main() {
	serverEnabledFlag := flag.Bool("server", false, "for running a server")
	var socketAddr string
	flag.StringVar(&socketAddr, "S", sockAddrDefault, "the unix socket to bind with")
	flag.Parse()
	if *serverEnabledFlag {
		s := NewServer()
		s.SocketAddr = socketAddr
		s.Listen()
	} else {
		s := NewSender(socketAddr)
		reader := s.inputChannel()
		data, _ := reader.ReadBytes('\n')
		fmt.Println(string(data))
		s.SendMessage(data)
	}

}
