package structs

import (
	"fmt"
	"log"
	"os/user"
)

func getUser() *user.User {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr
}

var User = getUser()

var SockAddrDefault = fmt.Sprintf(
	"%s/klystron.sock",
	User.HomeDir)

const (
	SockBufferSize64 = 212992
	SockBufferSize32 = 163840
	JobSpoolDir      = "/var/tmp/klystron"
)

type FontConfig struct {
	Name    string `json:"name"`
	Style   string `json:"style"`
	Size    int    `json:"size"`
	FontDir string `json:"fontDir"`
}

type PdfFile struct {
	Orientation string     `json:"orientation"`
	Unit        string     `json:"unit"`
	Format      string     `json:"format"`
	Font        FontConfig `json:"font"`
}

type DataSource struct {
	Path string `json:"path"`
	Data string `json:"data"`
	Hash string `json:"hash"`
}

type Job struct {
	Name      string       `json:"name"`
	Directory string       `json:"directory"`
	Outputs   []PdfFile    `json:"outputs"`
	Sources   []DataSource `json:"sources"`
}

type Message struct {
	Content Job    `json:"content"`
	Hash    string `json:"hash"`
}
