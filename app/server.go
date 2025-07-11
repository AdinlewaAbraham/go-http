package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/http-server-starter-go/internals/constants"
	"github.com/codecrafters-io/http-server-starter-go/internals/httpx"
	"github.com/codecrafters-io/http-server-starter-go/internals/server"
)

func main() {
	dir := flag.String("directory", "", "Path to the directory containing files")
	flag.Parse()
	fmt.Println("Logs from your program will appear here!")

	httpServer := server.CreateServer()

	httpServer.Get("/", func(req *httpx.HttpRequest, res *httpx.HttpResponse) {
		res.Status(200).Send([]byte("cool"))
	})

	httpServer.Get("/ping", func(req *httpx.HttpRequest, res *httpx.HttpResponse) {
		buff := bytes.Buffer{}

		for k, v := range req.Header {
			buff.WriteString(k)
			buff.WriteString(": ")
			buff.WriteString(v)
			buff.WriteString("\r\n")
		}

		res.Send(buff.Bytes())
	})

	httpServer.Get("/echo/:id", func(req *httpx.HttpRequest, res *httpx.HttpResponse) {
		id, exists := req.Params["id"]
		if !exists {
			res.Status(500).Send([]byte("Could not find param 'id'"))
			return
		}

		res.Status(200).SetHeader(string(constants.HeaderKeyContentType), string(constants.ContentTypeText))
		res.Send([]byte(id))
	})

	httpServer.Get("/user-agent", func(req *httpx.HttpRequest, res *httpx.HttpResponse) {
		res.SetHeader("Content-Type", "text/plain")
		res.Status(200).Send([]byte(req.UserAgent))
	})

	httpServer.Get("/files/:filename", func(req *httpx.HttpRequest, res *httpx.HttpResponse) {
		if *dir == "" {
			fmt.Println("Error: --directory flag is required")
			return
		}

		if _, exists := req.Params["filename"]; !exists {
			res.Status(500).Send([]byte("Could not find param 'filename'"))
			return
		}

		fullPath := filepath.Join(*dir, req.Params["filename"])
		fileData, err := os.ReadFile(fullPath)

		if err != nil {
			if os.IsNotExist(err) {
				res.Status(int(constants.StatusNotFound)).Send([]byte("File not found"))
			} else {
				res.Status(int(constants.StatusInternalServerError)).Send([]byte("Internal server error"))
			}
			return
		}

		res.SetHeader(string(constants.HeaderKeyContentType), string(constants.ContentTypeOctet))
		res.SendFile(req.Params["filename"], fileData)
	})

	httpServer.Post("/files/:filename", func(req *httpx.HttpRequest, res *httpx.HttpResponse) {
		if *dir == "" {
			fmt.Println("Error: --directory flag is required")
			return
		}

		if _, exists := req.Params["filename"]; !exists {
			fmt.Println("Could not find param 'filename")
			log.Printf("Could not find param filename, params map: %+v\n", req.Params)
			res.Status(500).Send([]byte("Could not find param 'filename'"))
			return
		}

		fullPath := filepath.Join(*dir, req.Params["filename"])
		file, err := os.Create(fullPath)

		if err != nil {
			log.Println("Could not create file: ", err.Error())
			res.Status(int(constants.StatusInternalServerError)).Send([]byte("Internal server error"))
			return
		}

		log.Printf("Raw req.Body: %q", req.Body)
		n, err := file.Write(req.Body)
		if err != nil {
			log.Printf("Error writing to file: %v", err)
			res.Status(500).Send([]byte("Internal Server Error"))
			return
		}
		file.Close()
		log.Printf("Wrote %s of size %d bytes to file", req.Body, n)

		res.Status(201).End()
	})

	httpServer.Listen("0.0.0.0:4221", func() {
		fmt.Println("listen callback")
	})
}
