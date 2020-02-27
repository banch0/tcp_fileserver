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
	client, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		t.Log("can't connect to the server: ", err)
	}
	defer client.Close()
	wg.Add(1)
	err = showDirectory(client, wg)
	if err != nil {
		if err == io.EOF {
			return
		}
		t.Log("clientDownload showDirectory Error: ", err)
	}

}

func TestDownload(t *testing.T) {
	client, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		t.Log("can't connect to the server: ", err)
	}
	defer client.Close()

	go downloadingFiles(client, "image.jpg")
	time.Sleep(100 * time.Millisecond)
}

func TestUpload(t *testing.T) {
	client, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		t.Log("can't connect to the server: ", err)
	}
	defer client.Close()

	go uploadingFiles(client, "image.jpg")
	time.Sleep(100 * time.Millisecond)
}
