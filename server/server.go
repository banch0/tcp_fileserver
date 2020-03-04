package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

// FileObject ...
type FileObject struct {
	FileName   string
	ReadWriter io.ReadWriter
}

const (
	// IP local addres
	IP = "0.0.0.0"
	// PORT number of server
	PORT = "9999"
	// UPLOAD ...
	UPLOAD = "upload"
	// DOWNLOAD ...
	DOWNLOAD = "download"
	// SHOWDIR ...
	SHOWDIR = "showdir"
)

func main() {
	listen, err := net.Listen("tcp", IP+":"+PORT)
	if err != nil {
		log.Fatal("Server can't started", err)
	}

	log.Printf("Server starting on %s\n", IP+":"+PORT)
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		log.Println("client connected")
		if err != nil {
			log.Println("Can't connect to the server")
		}
		go handleRequest(conn)
	}
}

// handleRequest ...
func handleRequest(conn net.Conn) {
	defer conn.Close()

	for {
		reader := bufio.NewReader(conn)
		command, err := reader.ReadString(' ')
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Println("Error read space:", err)
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("Error when reading client data %v", err)
			return
		}

		switch strings.TrimSpace(command) {
		case UPLOAD:
			log.Println("uploading ...")
			uploadFile(conn, line)
			return
		case DOWNLOAD:
			log.Println("downloading ...")
			downloadFile(conn, line)
			return
		case SHOWDIR:
			log.Println("showing dir ...")
			showDirectoryOnServer(conn)
			return
		}
	}

}

func downloadFile(conn net.Conn, line string) {
	file, err := os.Open("../files/" + strings.TrimSpace(line))
	if err != nil {
		log.Println("downloadFile os.Open Error: ", err)
	}
	defer file.Close()

	m := &FileObject{
		ReadWriter: file,
		FileName:   line,
	}

	n, err := io.Copy(conn, m.ReadWriter)
	if err != nil {
		log.Println("downlaodFile io.Copy Error:", err)
	}
	
	log.Println("Send bytes: ", n)
}

func uploadFile(conn net.Conn, line string) {
	file, err := os.Create("../files/" + strings.TrimSpace(line))
	if err != nil {
		log.Println("uploadFile os.Create Error: ", err)
	}
	defer file.Close()

	m := &FileObject{
		ReadWriter: file,
		FileName:   line,
	}

	n, err := io.Copy(m.ReadWriter, conn)
	if err != nil {
		log.Println("uploadFile io.Copy Error:", err)
	}
	log.Println("Download bytes:", n)
}

func showDirectoryOnServer(conn net.Conn) {
	var mystr string
	files, err := ioutil.ReadDir("../files/")

	if err != nil {
		log.Fatal("Can't open directory ", err)
	}

	for _, f := range files {
		mystr += " " + f.Name()
	}

	conn.Write([]byte(mystr))
}
