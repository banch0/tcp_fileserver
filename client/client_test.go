package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"sync"
	"testing"
)

const (
	ADDR = "127.0.0.1"
)

var serverFileSize int64

func TestShowdir(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		var mystr string
		listen, err := net.Listen("tcp", ADDR+":"+port)
		if err != nil {
			log.Println("creatServer error: ", err)
		}
		defer func() {
			err := listen.Close()
			if err != nil {
				t.Error("Close error", err)
			}
		}()

		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Println("Can't accept client error: ", err)
			}

			files, err := ioutil.ReadDir("../files/")
			if err != nil {
				log.Fatal("Can't open directory ", err)
			}

			for _, f := range files {
				mystr += " " + f.Name()
			}

			conn.Write([]byte(mystr))
		}
	}()

	go func() {
		client, err := net.Dial("tcp", address+":"+port)
		if err != nil {
			t.Log("can't connect to the server: ", err)
		}

		defer func() {
			err := client.Close()
			if err != nil {
				t.Error("Close error", err)
			}
		}()

		datas := make([]byte, 512)
		_, err = client.Read(datas)
		if err != nil {
			if err == io.EOF {
				return
			}
			t.Error("Error client read: ", err)
		}

		filesName := strings.Split(string(datas[1:]), " ")
		if len(filesName) != 2 {
			t.Error("Client showdir test error")
		}
		defer wg.Done()
	}()

	wg.Wait()
	// t.Error("end of showdir test")
}

func TestDownload(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		listen, err := net.Listen("tcp", ADDR+":"+port)
		if err != nil {
			log.Println("creatServer error: ", err)
		}

		defer func() {
			err := listen.Close()
			if err != nil {
				t.Error("Close error", err)
			}
		}()

		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Println("Can't accept client error: ", err)
			}

			reader := bufio.NewReader(conn)
			_, err = reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					return
				}
				log.Println("Error read space:", err)
			}

			// file, err := os.Open("../files/image.jpg")
			// if err != nil {
			// 	log.Println("downloadFile os.Open Error: ", err)
			// }
			// defer file.Close()

			// serverFileSize, err = io.Copy(conn, file)
			// if err != nil {
			// 	t.Error("Can't copy file")
			// }
			// conn.Write([]byte(strconv.Itoa(int(serverFileSize))))
			t.Log(serverFileSize)
		}
	}()

	go func() {
		client, err := net.Dial("tcp", address+":"+port)
		if err != nil {
			t.Log("can't connect to the server: ", err)
		}
		defer func() {
			err := client.Close()
			if err != nil {
				t.Error("Close error", err)
			}
		}()

		// var clientFileSize int64

		// file, err := os.Create("./file")

		// m := FileObject{
		// 	Source: file,
		// }

		// clientFileSize, err = io.Copy(m.Source, client)
		// if err != nil {
		// 	t.Error("Client can't copy file")
		// }

		// if clientFileSize == serverFileSize {
		// 	t.Error("Wrong size of downloading file")
		// }

		client.Write([]byte("\n"))
		defer wg.Done()
	}()
	wg.Wait()
	// t.Error("end of download test")
}

func TestUpload(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		listen, err := net.Listen("tcp", ADDR+":"+port)
		if err != nil {
			t.Log("creatServer error: ", err)
		}

		defer func() {
			err := listen.Close()
			if err != nil {
				t.Error("Close error", err)
			}
		}()

		for {
			conn, err := listen.Accept()
			if err != nil {
				t.Log("Can't accept client error: ", err)
			}

			reader := bufio.NewReader(conn)
			_, err = reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					return
				}
				log.Println("Error read space:", err)
			}

			// go func() {
			// 	file, err := os.Create("./testFile")
			// 	if err != nil {
			// 		log.Println("uploadFile os.Create Error: ", err)
			// 	}
			// 	defer file.Close()
			// 	m := FileObject{
			// 		Source: file,
			// 	}
			// 	n, err := io.Copy(m.Source, conn)
			// 	if err != nil {
			// 		t.Error("Can't create error")
			// 	}
			// 	t.Error(n)
			// }()
		}
	}()
	go func() {
		client, err := net.Dial("tcp", address+":"+port)
		if err != nil {
			t.Log("can't connect to the server: ", err)
		}
		defer func() {
			err := client.Close()
			if err != nil {
				t.Error("Close error", err)
			}
		}()

		client.Write([]byte("\n"))
		defer wg.Done()
	}()
	wg.Wait()
	// t.Error("end of upload test")
}
