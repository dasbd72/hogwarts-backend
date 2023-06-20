package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	MAX_SERVER int = 5
	MAX_CLIENT int = 10
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	mu_client_conn_list sync.Mutex
	mu_server_conn_list sync.Mutex
	client_conn_list    []*websocket.Conn
	client_conn_mu_list []*sync.Mutex
	client_conn_ok_list []bool
	server_conn_list    []*websocket.Conn
	server_conn_mu_list []*sync.Mutex
	server_conn_ok_list []bool
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/server", server)
	http.HandleFunc("/client", client)

	go http.ListenAndServeTLS(":https", "go-server.crt", "go-server.key", nil)
	log.Fatal(http.ListenAndServe(":http", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello, World!")
}

func server(w http.ResponseWriter, r *http.Request) {
	var conn *websocket.Conn
	var err error

	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	log.Printf("[server] %s connected", r.RemoteAddr)

	mu_server_conn_list.Lock()
	if len(server_conn_list)-(MAX_SERVER-1) < 0 {
		server_conn_list = append(server_conn_list[0:], conn)
		server_conn_mu_list = append(server_conn_mu_list[0:], &sync.Mutex{})
		server_conn_ok_list = append(server_conn_ok_list[0:], true)
	} else {
		server_conn_list = append(server_conn_list[len(server_conn_list)-(MAX_SERVER-1):], conn)
		server_conn_mu_list = append(server_conn_mu_list[len(server_conn_mu_list)-(MAX_SERVER-1):], &sync.Mutex{})
		server_conn_ok_list = append(server_conn_ok_list[len(server_conn_ok_list)-(MAX_SERVER-1):], true)
	}
	mu_server_conn_list.Unlock()

	for {
		mt, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[server][error] %s\n", err)
			return
		}

		mu_client_conn_list.Lock()
		for i, client_conn := range client_conn_list {
			if !client_conn_ok_list[i] {
				continue
			}
			client_conn_mu_list[i].Lock()
			err = client_conn.WriteMessage(mt, data)
			client_conn_mu_list[i].Unlock()
			if err != nil {
				log.Printf("[client][error] %s\n", err)
				client_conn_ok_list[i] = false
			}
		}
		mu_client_conn_list.Unlock()
	}
}

func client(w http.ResponseWriter, r *http.Request) {
	var conn *websocket.Conn
	var err error

	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	log.Printf("[client] %s connected", r.RemoteAddr)

	mu_client_conn_list.Lock()
	if len(client_conn_list)-(MAX_CLIENT-1) < 0 {
		client_conn_list = append(client_conn_list[0:], conn)
		client_conn_ok_list = append(client_conn_ok_list[0:], true)
		client_conn_mu_list = append(client_conn_mu_list[0:], &sync.Mutex{})
	} else {
		client_conn_list = append(client_conn_list[len(client_conn_list)-(MAX_CLIENT-1):], conn)
		client_conn_ok_list = append(client_conn_ok_list[len(client_conn_ok_list)-(MAX_CLIENT-1):], true)
		client_conn_mu_list = append(client_conn_mu_list[len(client_conn_mu_list)-(MAX_CLIENT-1):], &sync.Mutex{})
	}
	mu_client_conn_list.Unlock()

	for {
		mt, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[client][error] %s\n", err)
			return
		}

		mu_server_conn_list.Lock()
		for i, server_conn := range server_conn_list {
			if !server_conn_ok_list[i] {
				continue
			}
			server_conn_mu_list[i].Lock()
			err = server_conn.WriteMessage(mt, data)
			server_conn_mu_list[i].Unlock()
			if err != nil {
				log.Printf("[client][error] %s\n", err)
				server_conn_ok_list[i] = false
			}
		}
		mu_server_conn_list.Unlock()
	}
}
