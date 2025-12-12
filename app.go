// Package main provides a simple TCP server that responds to WHO queries with system information.
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
)

func main() {
	var (
		port              string
		certFile, keyFile string
		name              string
		banner            bool
	)

	flag.StringVar(&port, "port", ":8080", "give me a port number")
	flag.StringVar(&certFile, "certFile", "", "TLS - certificate path")
	flag.StringVar(&keyFile, "keyFile", "", "TLS - key path")
	flag.StringVar(&name, "name", "", "name")
	flag.BoolVar(&banner, "banner", false, "Connection banner")
	flag.Parse()

	log.Printf("Starting up on port %s", port)

	var listener net.Listener
	if certFile != "" && keyFile != "" {
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
		if banner {
			_, err := conn.Write([]byte("Welcome"))
			if err != nil {
				log.Fatal(err)
			}
		}
		go serveTCP(conn, name)
	}
}

func serveTCP(conn io.ReadWriteCloser, name string) {
	defer func() { _ = conn.Close() }()

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
			_, err := conn.Write([]byte(whoAmIInfo(name)))
			if err != nil {
				log.Println(err)
			}
		} else {
			_, err := fmt.Fprintf(conn, "Received: %s", buffer[:n])
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func whoAmIInfo(name string) string {
	var out bytes.Buffer

	if name != "" {
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
	config := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	config.Certificates = make([]tls.Certificate, 1)

	var err error
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)

	return config, err
}
