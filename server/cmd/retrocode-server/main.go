package main

import (
	"encoding/json"
	"io"
	"log"
	"net"

	"retrocode/server/internal/rpc"
)

const listenAddr = "127.0.0.1:7741"

func main() {
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("retrocode-server M0 listening on %s", listenAddr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	log.Printf("client connected: %s", conn.RemoteAddr())
	header := make([]byte, 4)
	for {
		if _, err := io.ReadFull(conn, header); err != nil {
			return
		}
		length, compressed, err := rpc.DecodeHeader(header)
		if err != nil {
			log.Printf("bad header: %v", err)
			return
		}
		payload := make([]byte, length)
		if _, err := io.ReadFull(conn, payload); err != nil {
			return
		}
		var msg map[string]interface{}
		if err := json.Unmarshal(payload, &msg); err != nil {
			log.Printf("bad json: %v", err)
			return
		}
		log.Printf("recv method=%v compressed=%v bytes=%d", msg["m"], compressed, length)
		resp, _ := json.Marshal(map[string]interface{}{
			"m": "error",
			"params": map[string]interface{}{
				"code":    "not_implemented",
				"message": "M0: transport up, methods land in M1",
			},
		})
		frame, _ := rpc.Encode(resp, false)
		if _, err := conn.Write(frame); err != nil {
			return
		}
	}
}
