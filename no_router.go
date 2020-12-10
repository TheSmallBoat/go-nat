package nat

import (
	"fmt"
	"net"
)

const googleDNSServer = "8.8.8.8:80"

func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", googleDNSServer)
	if err != nil {
		return nil, err
	}

	if udpAddr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
		return udpAddr.IP, conn.Close()
	}

	_ = conn.Close()
	return nil, fmt.Errorf("getting outbound IP failed")
}
