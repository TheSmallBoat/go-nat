package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/thesmallboat/go-nat"
)

func main() {
	n, err := nat.DiscoverGateway()
	if err != nil {
		ip, err := nat.GetOutboundIP()
		if err != nil {
			log.Printf("%v", err)
		} else {
			log.Printf("getting outbound IP: %v", ip)
		}

		log.Fatalf("error: %s", err)
	}
	log.Printf("nat type: %s", n.Type())

	daddr, err := n.GetDeviceAddress()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	log.Printf("device address: %s", daddr)

	iaddr, err := n.GetInternalAddress()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	log.Printf("internal address: %s", iaddr)

	eaddr, err := n.GetExternalAddress()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	log.Printf("external address: %s", eaddr)

	eport, err := n.AddPortMapping("tcp", 8888, "http", 60)
	if err != nil {
		// Some hardware does not support mappings with timeout, so try that
		eport, err = n.AddPortMapping("tcp", 8888, "http", 0)
		if err != nil {
			log.Fatalf("error: %s", err)
		}
	}

	log.Printf("test-page: http://%s:%d/", eaddr, eport)

	go func() {
		for {
			time.Sleep(30 * time.Second)

			_, err = n.AddPortMapping("tcp", 8888, "http", 60)
			if err != nil {
				log.Fatalf("error: %s", err)
			}
		}
	}()

	defer n.DeletePortMapping("tcp", 8888)

	http.ListenAndServe(":8888", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(200)
		fmt.Fprintf(rw, "Hello there!\n")
		fmt.Fprintf(rw, "nat type: %s\n", n.Type())
		fmt.Fprintf(rw, "device address: %s\n", daddr)
		fmt.Fprintf(rw, "internal address: %s\n", iaddr)
		fmt.Fprintf(rw, "external address: %s\n", eaddr)
		fmt.Fprintf(rw, "test-page: http://%s:%d/\n", eaddr, eport)
	}))
}
