package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/himanshuraimau/drift/p2p"
)

// makeServer creates a new file server instance
func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		EncKey:            newEncryptionKey(),
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {
	// Create three nodes
	s1 := makeServer(":3000", "")
	s2 := makeServer(":7000", "")
	s3 := makeServer(":5000", ":3000", ":7000")

	// Start the first two nodes
	go func() {
		log.Fatal(s1.Start())
	}()
	time.Sleep(500 * time.Millisecond)

	go func() {
		log.Fatal(s2.Start())
	}()
	time.Sleep(2 * time.Second)

	// Start the third node and let it connect to the network
	go s3.Start()
	time.Sleep(2 * time.Second)

	// Demonstrate file operations
	fmt.Println("=== Distributed File System Demo ===")

	// Store files
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("picture_%d.png", i)
		data := bytes.NewReader([]byte(fmt.Sprintf("This is the content of file %d - distributed file system data!", i)))
		
		fmt.Printf("Storing file: %s\n", key)
		if err := s3.Store(key, data); err != nil {
			log.Printf("Error storing file %s: %v\n", key, err)
			continue
		}

		// Retrieve the file
		fmt.Printf("Retrieving file: %s\n", key)
		r, err := s3.Get(key)
		if err != nil {
			log.Printf("Error retrieving file %s: %v\n", key, err)
			continue
		}

		b, err := io.ReadAll(r)
		if err != nil {
			log.Printf("Error reading file %s: %v\n", key, err)
			continue
		}

		fmt.Printf("File content: %s\n", string(b))
		
		// Close the reader if it's a ReadCloser
		if rc, ok := r.(io.ReadCloser); ok {
			rc.Close()
		}

		// Delete the file
		fmt.Printf("Deleting file: %s\n", key)
		if err := s3.Delete(key); err != nil {
			log.Printf("Error deleting file %s: %v\n", key, err)
		}

		fmt.Println("---")
	}

	fmt.Println("=== Demo Complete ===")
	
	// Keep the servers running for a while
	time.Sleep(5 * time.Second)
}
