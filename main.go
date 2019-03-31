// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var addr = flag.String("addr", "0.0.0.0:80", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	fmt.Println("running...")
	flag.Parse()

	strConn := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s",
		os.Getenv("C2G_DB_USER"),
		os.Getenv("C2G_DB_NAME"),
		os.Getenv("C2G_DB_PASSWORD"),
		os.Getenv("C2G_DB_HOST"),
		os.Getenv("C2G_DB_PORT"),
	)

	db, err := sql.Open("postgres", strConn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	hub := newHub(db)

	go hub.run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
