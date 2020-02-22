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

/*
// TODO:
1. Client.exe upload <paht_to_file> +
2. Server.exe download <path_to_file> +
3. Client.exe list ("show list of files") +

// needed packages is: filepath -, flag+
// Warning: all requests must run in goroutines
// test run from goroutines too
// add to github and add Actions
*/

// MyConnection ...
type MyConnection struct {
	Num        int64
	connection *net.Conn
	FileWriter io.Writer
	FileReader io.Reader
	Datas      []byte
	File       *os.File
}

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
var PORT = ":9999"

func main() {
	listen, err := net.Listen("tcp", IP+PORT)
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
	file, err := os.Open(strings.TrimSpace(line))
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
	file, err := os.Create("./" + strings.TrimSpace(line))
	if err != nil {
		log.Println("uploadFile os.Create Error: ", err)
	}
	defer file.Close()
	m := &MyFile{ReadWriter: file, FileName: line, }
	n, err := io.Copy(m.ReadWriter, conn)
	if err != nil {
		log.Println("uploadFile io.Copy Error:", err)
	}
	log.Println("Download bytes:", n)
	return
}

func showDirectory(conn net.Conn) {
	var mystr string
	files, err := ioutil.ReadDir("./")

	if err != nil {
		log.Fatal("Can't open directory ", err)
	}

	for _, f := range files {
		mystr += " " + f.Name()
	}

	conn.Write([]byte(mystr))
	return
}
