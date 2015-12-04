package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sort"
)


func ipHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(strings.Split(r.RemoteAddr, ":")[0]))
}

func portHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(strings.Split(r.RemoteAddr, ":")[1]))
}

func ipPortHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.RemoteAddr))
}

func headersHandler(w http.ResponseWriter, r *http.Request) {

	// sort headers
	headers := make([]string, len(r.Header))
	i := 0
	for k, _ := range r.Header {
		headers[i] = k
		i++
    }
	sort.Strings(headers)

	// print headers
	for _, h := range headers {
		fmt.Fprintf(w, "%s: %s\n", h, strings.Join(r.Header[h], "; "))
	}
}

func main() {

	addr := ":8181";
	if (len(os.Args) > 1) {
		addr = os.Args[1]
	}

	http.HandleFunc("/", ipHandler)
	http.HandleFunc("/port", portHandler)
	http.HandleFunc("/ipport", ipPortHandler)
	http.HandleFunc("/headers", headersHandler)

	err := http.ListenAndServe(addr, nil)

	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
