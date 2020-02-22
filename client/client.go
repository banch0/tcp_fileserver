package main

import (
	"bufio"
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	// "time"

	"io"
	"log"
	"net"
	"os"
)

var upload string
var download string
var showList bool

// IP addres of server
var IP = "127.0.0.1"

// PORT number of server
var PORT = ":9999"

func init() {
	flag.StringVar(&upload, "upload", "", "Add path to file")
	flag.StringVar(&upload, "u", "", "Add path to file")
	flag.StringVar(&download, "download", "", "Add path to file")
	flag.StringVar(&download, "d", "", "Add path to file")
	flag.BoolVar(&showList, "list", false, "show all file in directory")
	flag.BoolVar(&showList, "l", false, "show all file in directory")
}

// MyFile ...
type MyFile struct {
	FileName string
	Source   io.Writer
	Files    io.Reader
	Length   int64
}

func main() {
	flag.Parse()
	wg := &sync.WaitGroup{}
	client, err := net.Dial("tcp", IP+PORT)
	if err != nil {
		log.Fatal("Can't connect to the server", err)
	}

	// defer client.Close()
	if showList != false {
		wg.Add(1)
		// time.Sleep(5 * time.Second)
		// go showDirectory(client, wg)
		go clientDownload(wg)
		log.Println("in the show")
		// return
	}

	wg.Wait()

	if upload != "" {
		go uploadingFiles(client, upload)
		time.Sleep(time.Second * 2)
		return
	}

	if download != "" {
		go downloadingFiles(client, download)
		time.Sleep(time.Second * 2)
		return
	}
}

func clientDownload(wg *sync.WaitGroup) {
	client, err := net.Dial("tcp", IP+PORT)
	if err != nil {
		log.Fatal("Can't connect to the server", err)
		return
	}

	showDirectory(client, wg)

	defer client.Close()
	// return err
}

func showDirectory(client net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	client.Write([]byte("showdir \n"))
	d := make([]byte, 512)
	_, err := client.Read(d)
	if err != nil {
		if err == io.EOF {
			log.Println("endOFfile")
			return
		}
		log.Println("Error client read", err)
	}
	fl := strings.Split(string(d[1:]), " ")
	fmt.Println(" Files on the server:")

	for _, file := range fl {
		fmt.Printf(" - %s\n", file)
	}
	return
}

func downloadingFiles(client net.Conn, download string) {
	log.Println("start downloading")

	out, err := os.Create("./" + download)
	if err != nil {
		log.Println("download Create Error: ", err)
	}
	defer out.Close()

	writer := bufio.NewWriter(client)

	writer.WriteString("download " + download + "\n")
	writer.Flush()

	m := &MyFile{Source: out, FileName: download}

	n, err := io.Copy(m.Source, client)
	if err != nil {
		if err == io.EOF {
			return
		}
		log.Println("download Error: ", err)
	}
	fmt.Println("Recieve bytes: ", n)
}

func uploadingFiles(client net.Conn, upload string) {
	log.Println("start uploading")

	file, err := os.Open(strings.TrimSpace(upload))
	if err != nil {
		log.Println("Open file error:", err)
	}
	defer file.Close()

	client.Write([]byte("upload " + upload + "\n"))
	time.Sleep(5 * time.Millisecond)

	n, err := io.Copy(client, file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Send bytes:", n)
}
