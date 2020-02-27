package main

import (
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
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
			conn.Write([]byte("echo"))
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
		err = showDirectory(client, wg)
		if err != nil {
			if err == io.EOF {
				return
			}
			t.Log("clientDownload showDirectory Error: ", err)
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

		go uploadingFiles(client, "image.jpg")
		time.Sleep(100 * time.Millisecond)
	}()
}
