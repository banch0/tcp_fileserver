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
	"time"
)

func TestShowdir(t *testing.T) {
	var mystr string
	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		listen, err := net.Listen("tcp", address+":"+port)
		if err != nil {
			t.Log("creatServer error: ", err)
		}
		defer listen.Close()
		defer wg.Done()
		for {
			conn, err := listen.Accept()
			if err != nil {
				t.Log("Can't accept client error: ", err)
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
				t.Log("Close error", err)
			}
		}()
		wg.Add(1)
		client.Write([]byte("showdir \n"))
		datas := make([]byte, 512)
		_, err = client.Read(datas)
		if err != nil {
			if err == io.EOF {
				return
			}
			t.Log("Error client read: ", err)
		}
		filesName := strings.Split(string(datas[1:]), " ")
		if len(filesName) != 2 {
			t.Error("Client showdir test error")
		}
		defer wg.Done()
	}()
	wg.Wait()

}

func TestDownload(t *testing.T) {
	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		listen, err := net.Listen("tcp", address+":"+port)
		if err != nil {
			t.Log("creatServer error: ", err)
		}
		defer listen.Close()
		defer wg.Done()
		for {
			_, err := listen.Accept()
			t.Log("server listen")
			if err != nil {
				t.Log("Can't accept client error: ", err)
			}
		}
	}()
	go func() {
		client, err := net.Dial("tcp", address+":"+port)
		if err != nil {
			t.Log("can't connect to the server: ", err)
		}
		defer client.Close()

		go downloadingFiles(client, "image.jpg")
		time.Sleep(100 * time.Millisecond)
	}()
}

func TestUpload(t *testing.T) {
	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		listen, err := net.Listen("tcp", address+":"+port)
		if err != nil {
			t.Log("creatServer error: ", err)
		}
		defer listen.Close()
		defer wg.Done()
		for {
			conn, err := listen.Accept()
			if err != nil {
				t.Log("Can't accept client error: ", err)
			}
			reader := bufio.NewReader(conn)
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					return
				}
				t.Logf("Error when reading client data %v", err)
				return
			}
			t.Log(line)
		}
	}()
	go func() {
		client, err := net.Dial("tcp", address+":"+port)
		if err != nil {
			t.Log("can't connect to the server: ", err)
		}
		defer client.Close()

		go uploadingFiles(client, "image.jpg")
		time.Sleep(100 * time.Millisecond)
	}()
}
