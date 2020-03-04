package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	showList bool
	download string
	upload   string
	address  = "localhost"
	port     = "9999"
)

// SHOWDIR ...
const SHOWDIR = "showdir"

func init() {
	flag.StringVar(&upload, "upload", "", "Add path to file")
	flag.StringVar(&upload, "u", "", "Add path to file")
	flag.StringVar(&download, "download", "", "Add path to file")
	flag.StringVar(&download, "d", "", "Add path to file")
	flag.BoolVar(&showList, "list", false, "show all file in directory")
	flag.BoolVar(&showList, "l", false, "show all file in directory")
}

// FileObject ...
type FileObject struct {
	FileName string
	Source   io.ReadWriter
}

func main() {
	flag.Parse()
	wg := &sync.WaitGroup{}

	client, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		log.Fatal("Can't connect to the server", err)
	}
	defer client.Close()

	if showList != false {
		wg.Add(1)
		go showDirectoryClient(wg)
	}

	if upload != "" {
		go uploadingFiles(client, upload)
		time.Sleep(time.Millisecond * 100)
		return
	}

	if download != "" {
		go downloadingFiles(client, download)
		time.Sleep(time.Millisecond * 100)
		return
	}

	wg.Wait()
}

func showDirectoryClient(wg *sync.WaitGroup) error {
	client, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		log.Fatal("Can't connect to the server", err)
		return err
	}
	defer client.Close()

	err = showDirectory(client, wg)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		log.Println("showDirectoryClient showDirectory Error: ", err)
	}
	return err
}

func showDirectory(client net.Conn, wg *sync.WaitGroup) error {
	defer wg.Done()

	client.Write([]byte(SHOWDIR + "\n"))
	datas := make([]byte, 512)

	_, err := client.Read(datas)
	if err != nil {
		if err == io.EOF {
			log.Println("showDirectory end Of file")
			return err
		}
		log.Println("Error client read", err)
	}

	filesName := strings.Split(string(datas[1:]), " ")
	fmt.Println(" Files on the server:")

	for _, file := range filesName {
		fmt.Printf(" - %s\n", file)
	}
	
	return err
}

func downloadingFiles(client net.Conn, download string) {
	fmt.Println("--- Start downloading")

	out, err := os.Create("./" + download)

	if err != nil {
		log.Println("download Create Error: ", err)
	}

	defer out.Close()

	writer := bufio.NewWriter(client)

	writer.WriteString("download " + download + "\n")
	writer.Flush()

	m := &FileObject{
		Source:   out,
		FileName: download,
	}

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
	fmt.Println("--- Start uploading")

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
