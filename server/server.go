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

// MyFile ...
type MyFile struct {
	FileName   string
	Source     io.Writer
	Files      io.Reader
	ReadWriter io.ReadWriter
	Length     int64
	Datas      []byte
}

// IP addres of server
var IP = "0.0.0.0"

// PORT number of server
var PORT = "9999"

func main() {
	listen, err := net.Listen("tcp", IP+":"+PORT)
	if err != nil {
		log.Fatal("Server can't started", err)
	}

	log.Printf("Server starting on %s\n", IP+PORT)
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
		case "upload":
			log.Println("uploading ...")
			uploadFile(conn, line)
			return
		case "download":
			log.Println("downloading ...")
			downloadFile(conn, line)
			return
		case "showdir":
			log.Println("showing dir ...")
			showDirectory(conn)
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
	m := &MyFile{Files: file, FileName: line}
	n, err := io.Copy(conn, m.Files)
	if err != nil {
		log.Println("downlaodFile io.Copy Error:", err)
	}
	log.Println("Send bytes: ", n)
	return
}

func uploadFile(conn net.Conn, line string) {
	file, err := os.Create("../files/" + strings.TrimSpace(line))
	if err != nil {
		log.Println("uploadFile os.Create Error: ", err)
	}
	defer file.Close()
	m := &MyFile{ReadWriter: file, FileName: line}
	n, err := io.Copy(m.ReadWriter, conn)
	if err != nil {
		log.Println("uploadFile io.Copy Error:", err)
	}
	log.Println("Download bytes:", n)
	return
}

func showDirectory(conn net.Conn) {
	var mystr string
	files, err := ioutil.ReadDir("../files/")

	if err != nil {
		log.Fatal("Can't open directory ", err)
	}

	for _, f := range files {
		mystr += " " + f.Name()
	}

	conn.Write([]byte(mystr))
	return
}
