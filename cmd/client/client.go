package main

import (
	_ "bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	/* net
	   	conn, err := net.Dial("tcp", "localhost:8080")
	       if err != nil {
	       	log.Fatalln(err)
	       }
	       fmt.Fprintf(conn, "GET /run HTTP/1.0\r\n\r\n")
	       status, err := bufio.NewReader(conn).ReadString('\n')
	   	fmt.Println(status)
	*/

	//http
	/*
	   	resp, err := http.Get("123123123")
	   	if err != nil {
	   		log.Fatalln(err)
	   	}
	   	defer resp.Body.Close()

	       body, err := io.ReadAll(resp.Body)
	   	if err != nil {
	   		log.Fatalln(err)
	   	}

	   	fmt.Println(string(body), resp.Header)
	   	fmt.Println(resp.Status)
	*/

	s := &http.Server{
		Addr:           ":8000",
		Handler:        NewRouter(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/test1", handler1)
	router.HandleFunc("/test2", handler2)
	router.HandleFunc("/test3", handler3)
	router.HandleFunc("/test4", handler4)
	router.HandleFunc("/test5", handler5)
	router.HandleFunc("/test6", handler6)

	return router
}

func handler1(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get req1")
	io.WriteString(w, "11111111111")
}

func handler2(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get req2")
	io.WriteString(w, "22222222222")
}

func handler3(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get req3")
	io.WriteString(w, "33333333333")
}

func handler4(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get req4")
	io.WriteString(w, "44444444444")
}

func handler5(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get req5")
	io.WriteString(w, "555555555555")
}

func handler6(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get req6")
	io.WriteString(w, "666666666666")
}
