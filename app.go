package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/mailgun/proxyproto"
)

var port string
var certFile, keyFile string
var name string
var banner, proxyprotocol bool

func init() {
	flag.StringVar(&port, "port", ":8080", "give me a port number")
	flag.StringVar(&certFile, "certFile", "", "TLS - certificate path")
	flag.StringVar(&keyFile, "keyFile", "", "TLS - key path")
	flag.StringVar(&name, "name", "", "name")
	flag.BoolVar(&banner, "banner", false, "Connection banner")
	flag.BoolVar(&proxyprotocol, "proxyprotocol", false, "Enable Proxy Protocol")
}

func main() {
	flag.Parse()
	fmt.Printf("Starting up on port %s", port)

	var listener net.Listener
	if len(certFile) > 0 && len(keyFile) > 0 {
		tlsConfig, err := createTLSConfig(certFile, keyFile)
		if err != nil {
			log.Fatalf("error creating TLS configuration: %v", err)
		}

		listener, err = tls.Listen("tcp", port, tlsConfig)
		if err != nil {
			log.Fatalf("error opening port: %v", err)
		}
	} else {
		var err error
		listener, err = net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("error opening port: %v", err)
		}
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go serveTCP(conn)
	}
}

func serveTCP(conn io.ReadWriteCloser) {
	defer conn.Close()
	var err error

	h := &proxyproto.Header{}
	if proxyprotocol {
		h, err = proxyproto.ReadHeader(conn)
		if err != nil {
			log.Printf("[ERR] Proxy protocol: %s", err)
			return
		}
	}

	if banner {
		_, err := fmt.Fprintln(conn, "Welcome")
		if err != nil {
			log.Fatal(err)
		}
	}

	for {
		buffer := make([]byte, 256)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println(err)
			return
		}

		temp := strings.TrimSpace(string(buffer[:n]))
		if temp == "STOP" {
			break
		}

		if temp == "WHO" {
			_, err := conn.Write([]byte(whoAmIInfo(h)))
			if err != nil {
				log.Println(err)
			}
		} else {
			_, err := conn.Write([]byte(fmt.Sprintf("Received: %s", buffer[:n])))
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func whoAmIInfo(h *proxyproto.Header) string {
	var out bytes.Buffer

	if h.Version > 0 {
		out.WriteString(fmt.Sprintf("Proxy Protocol Version: %d\n", h.Version))
		out.WriteString(fmt.Sprintf("Proxy Protocol Source: %s\n", h.Source.String()))
		out.WriteString(fmt.Sprintf("Proxy Protocol Destination: %s\n", h.Destination.String()))
		// add later: parse tlvs in proxy v2
		// if h.Version == 2 && len(h.RawTLVs) != 0 {
		//	tlvs, err := h.ParseTLVs()
		//	if err != nil {
		//       out.WriteString(fmt.Sprintf("Proxy Protocol TLVs: %d\n", ...))
		//	}
		//}
	}
	if len(name) > 0 {
		out.WriteString(fmt.Sprintf("Name: %s\n", name))
	}

	hostname, _ := os.Hostname()
	out.WriteString(fmt.Sprintf("Hostname: %s\n", hostname))

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			out.WriteString(fmt.Sprintf("IP: %s\n", ip))
		}
	}

	return out.String()
}

func createTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	config := &tls.Config{}
	config.Certificates = make([]tls.Certificate, 1)

	var err error
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)

	return config, err
}
