package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/thesmallboat/go-nat"
)

func main() {
	n, errA := nat.DiscoverGateway()
	if errA != nil {
		ip, errB := nat.GetOutboundIP()
		if errB != nil {
			log.Printf("error B: %v", errB)
		} else {
			log.Printf("getting outbound IP: %v", ip)
		}

		log.Fatalf("error A: %v", errA)
	}

	log.Printf("nat type: %v", n.Type())

	daddr, errC := n.GetDeviceAddress()
	if errC != nil {
		log.Fatalf("error C: %v", errC)
	}
	log.Printf("device address: %v", daddr)

	iaddr, errD := n.GetInternalAddress()
	if errD != nil {
		log.Fatalf("error D: %v", errD)
	}
	log.Printf("internal address: %v", iaddr)

	eaddr, errE := n.GetExternalAddress()
	if errE != nil {
		log.Fatalf("error E: %v", errE)
	}
	log.Printf("external address: %v", eaddr)

	eport, errF := n.AddPortMapping("tcp", 8888, "http", 60)
	if errF != nil {
		log.Fatalf("error F: %v", errF)
	}

	log.Printf("external port: %v", eport)
	log.Printf("test-page: http://%s:%d/", eaddr, eport)

	go func() {
		defer n.DeletePortMapping("tcp", 8888)

		for {
			time.Sleep(30 * time.Second)

			_, errG := n.AddPortMapping("tcp", 8888, "http", 60)
			if errG != nil {
				log.Fatalf("error G: %v", errG)
			}
		}
	}()

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
