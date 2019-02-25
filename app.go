package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var port string
var certFile, keyFile string
var name string

func init() {
	flag.StringVar(&port, "port", ":8080", "give me a port number")
	flag.StringVar(&certFile, "certFile", "", "TLS - certificate path")
	flag.StringVar(&keyFile, "keyFile", "", "TLS - key path")
	flag.StringVar(&name, "name", "", "name")
}

func main() {
	flag.Parse()

	fmt.Println("Starting up on port " + port)

	var listener net.Listener
	var err error
	if len(certFile) > 0 && len(keyFile) > 0 {
		tlsConfig, err := createTlsConfig(certFile, keyFile)
		if err != nil {
			log.Fatal("error creating TLS configuration: %v", err)
		}
		listener, err = tls.Listen("tcp", port, tlsConfig)
	} else {
		listener, err = net.Listen("tcp", port)
	}

	if err != nil {
		log.Fatal("error opening port: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go serveTCP(conn)
	}
}

func serveTCP(conn net.Conn) {
	defer conn.Close()

	for {
		bytes := make([]byte, 256)
		n, err := conn.Read(bytes)
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(string(bytes[:n]))
		if temp == "STOP" {
			break
		} else if temp == "WHO" {
			result := whoAmIInfo()
			conn.Write([]byte(result))
		} else {
			result := fmt.Sprintf("Received: %s", bytes[:n])
			conn.Write([]byte(result))
		}
	}
}

func whoAmIInfo () string {
	var out bytes.Buffer

	if len(name)>0 {
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

func createTlsConfig(certFile, keyFile string) (*tls.Config, error) {
	var err error

	config := &tls.Config{}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)

	return config, err
}
