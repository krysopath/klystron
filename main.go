package main

import (
	"flag"

	"github.com/krysopath/klystron/client"
	"github.com/krysopath/klystron/server"
	"github.com/krysopath/klystron/structs"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	serverEnabledFlag := flag.Bool(
		"server", false, "for running a server")
	jobFile := flag.String("job", "./Klystronfile", "a job file")

	var socketAddr string
	flag.StringVar(&socketAddr, "S",
		structs.SockAddrDefault, "the unix socket to bind with")

	flag.Parse()

	if *serverEnabledFlag {
		s := server.NewServer(socketAddr)
		s.Listen()
	} else {
		s := client.NewClient(socketAddr)
		s.Post(jobFile, flag.Args())
	}

}
