package client

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/krysopath/klystron/jobs"
	"github.com/krysopath/klystron/structs"
	"github.com/krysopath/klystron/utils"
	"gopkg.in/yaml.v2"
)

type Client struct {
	//Logger     *log.Logger
	SocketAddr string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (c *Client) sendMessage(message []byte) {
	conn, err := net.Dial("unix", c.SocketAddr)
	if err != nil {
		log.Printf("Failed to dial: %s", err)
		// return tcp/connection-refused status
		os.Exit(111)
	}
	defer conn.Close()
	count, err := conn.Write(message)
	if err != nil {
		log.Printf("Write error: %s", err)
	}
	log.Printf("Wrote %d bytes", count)
}

func (c *Client) inputChannel(files []string) []*bufio.Reader {
	var sources []*bufio.Reader
	for _, p := range files {
		fd, err := os.Open(p)
		check(err)
		sources = append(sources, bufio.NewReader(fd))
	}
	return sources
}

func (c *Client) Post(jobFile *string, files []string) {
	jobData, err := ioutil.ReadFile(*jobFile)
	check(err)
	job := jobs.JobUnmarshal(jobData)
	for i, reader := range c.inputChannel(files) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		data := buf.String()
		hashValue := utils.Hash(data)
		data = base64.StdEncoding.EncodeToString([]byte(data))
		//gzData := utils.GZipData([]byte(data))
		job.Sources = append(
			job.Sources,
			structs.DataSource{
				Path: files[i],
				Data: string(data),
				Hash: string(hashValue),
			},
		)
		log.Printf("encoding data source %d", i)
	}
	yml, err := yaml.Marshal(&job)
	check(err)
	fmt.Printf("---\n%s\n", string(yml))
	jobHash := utils.Hash(string(jobs.JobMarshal(&job)))
	message := structs.Message{
		Content: job,
		Hash:    jobHash,
	}
	c.sendMessage(utils.JSONMarshal(message))
}

func NewClient(socketAddr string) Client {
	return Client{SocketAddr: socketAddr}
}
