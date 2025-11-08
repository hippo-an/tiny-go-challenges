package main

import (
	"flag"
	"log"
	"os"

	"github.com/hippo-an/tiny-go-challenges/netpro_99/udp_reliable/tftp"
)

var (
	address = flag.String("a", "127.0.0.1:69", "listen address")
	payload = flag.String("p", "gopher.png", "file to serve to client")
)

func main() {
	flag.Parse()

	p, err := os.ReadFile(*payload)
	if err != nil {
		log.Fatal(err)
	}

	s := tftp.Server{Payload: p}
	log.Fatal(s.ListenAndServe(*address))

}
