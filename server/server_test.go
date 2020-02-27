package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

func TestShowdir(t *testing.T) {
	go func() {
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
	}()
	// wg := &sync.WaitGroup{}
	go func() {
		client, err := net.Dial("tcp", IP+":"+PORT)
		if err != nil {
			t.Error("can't connect to the server: ", err)
		}
		defer func() {
			err := client.Close()
			if err != nil {
				t.Error("Close error", err)
			}
		}()
		// wg.Add(1)

		// defer wg.Done()
		client.Write([]byte("showdir \n"))
		datas := make([]byte, 512)
		_, err = client.Read(datas)
		filesName := strings.Split(string(datas[1:]), " ")
		fmt.Println(" Files on the server:")

		for _, file := range filesName {
			fmt.Printf(" - %s\n", file)
		}
	}()
	// wg.Wait()
}

func TestDownloadFile(t *testing.T) {
	go func() {
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
	}()

	go func() {
		client, err := net.Dial("tcp", IP+":"+PORT)
		if err != nil {
			t.Error("can't connect to the server: ", err)
		}
		defer func() {
			err := client.Close()
			if err != nil {
				t.Error("Close error", err)
			}
		}()

		writer := bufio.NewWriter(client)
		writer.WriteString("download image.jpg\n")
		writer.Flush()
	}()
}

func TestUploadFile(t *testing.T) {
	go func() {
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
	}()

	go func() {
		client, err := net.Dial("tcp", IP+":"+PORT)
		if err != nil {
			t.Error("can't connect to the server: ", err)
		}
		defer func() {
			err := client.Close()
			if err != nil {
				t.Error("Close error", err)
			}
		}()
		file, err := os.Open("image.jpg")
		if err != nil {
			log.Println("Open file error:", err)
		}
		defer file.Close()
		client.Write([]byte("upload image.jpg\n"))
		time.Sleep(5 * time.Millisecond)

		n, err := io.Copy(client, file)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Send bytes:", n)
	}()
}
