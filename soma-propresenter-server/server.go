package main

import (
    "fmt"
    "log"
	"io"
	"regexp"
    "net/http"
	"github.com/gorilla/pat"
)

func macrosHandler(resWriter http.ResponseWriter, req *http.Request) {
	triggerMatch, err := regexp.MatchString("/macro/([a-fA-F0-9-]+)/trigger", req.URL.Path)
	fmt.Printf("matching: %s (%t)\n", req.URL.Path, triggerMatch)

	if err != nil {
		fmt.Printf("error in regex %s\n", err)
		panic(err)
	}
    if req.URL.Path != "/macros" && !triggerMatch {
        http.Error(resWriter, "404 not found.", http.StatusNotFound)
        return
    }

    if req.Method != "GET" {
        http.Error(resWriter, "Method is not supported.", http.StatusNotFound)
        return
    }

	requestURL := "http://localhost:1025/v1" + req.URL.Path
	fmt.Printf("Making http request: %s\n", requestURL)
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		panic(err)
	}
	resWriter.WriteHeader(res.StatusCode)
	io.Copy(resWriter, res.Body)
}

func main() {
	var mux *pat.Router = pat.New()
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/", fileServer)
	mux.Get("/macros", macrosHandler)
	mux.Get("/macro/{uuid}/trigger", macrosHandler)


	server := http.Server {
		Addr: ":8080",
		Handler: mux,
	}
    fmt.Printf("Starting server at port 8080\n")
    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}