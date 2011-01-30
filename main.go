package main

import (
	"http"
	"io"
	"log"
	// "fmt"
	"net"
)

func ProxyServer(w http.ResponseWriter, req *http.Request) {
	log.Printf("%#v", req)
	addr, err := net.ResolveIPAddr(req.URL.Host)
	if err != nil {
		io.WriteString(w, err.String())
		return
	}
	conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{ IP: addr.IP, Port: 80} )
	if err != nil {	
		io.WriteString(w, err.String())
		return
	}
	defer conn.Close()
	client := http.NewClientConn(conn, nil)
	err = client.Write(req)
	if err != nil {	
		io.WriteString(w, err.String())
		return
	}
	resp, err := client.Read()
	if err != nil {	
		io.WriteString(w, err.String())
		return
	}
	for k,v := range resp.Header {
		w.SetHeader(k, v)
	}
	io.Copy(w, resp.Body)
	w.Flush()
}
	

func main() {
	http.HandleFunc("/", ProxyServer)
	err := http.ListenAndServe(":8123", nil)
	if err != nil {
		log.Exit("ListenAndServe: ", err.String())
	}
}