package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sort"
)

var use_x_forwarded_for = false
var use_x_real_ip = false

func getRemoteIP(r *http.Request) string {

	if use_x_real_ip {
		realIP := r.Header.Get("X-Real-IP")
		if len(realIP) > 0 {
			return realIP
		}
	}

	if use_x_forwarded_for {
		forwardedFor := r.Header.Get("X-Forwarded-For")
		if len(forwardedFor) > 0 {
			return strings.Split(forwardedFor, ",")[0]
		}
	}

	return strings.Split(r.RemoteAddr, ":")[0]
}

func getPort(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[1]
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", getRemoteIP(r))
}

func portHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", getPort(r))
}

func ipPortHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s:%s\n", getRemoteIP(r), getPort(r))
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

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--use-x-real-ip":
			use_x_real_ip = true
		case "--use-x-forwarded-for":
			use_x_forwarded_for = true
		default:
			addr = arg;
		}
	}

	http.HandleFunc("/", ipHandler)
	http.HandleFunc("/port", portHandler)
	http.HandleFunc("/ipport", ipPortHandler)
	http.HandleFunc("/headers", headersHandler)

	err := http.ListenAndServe(addr, nil)

	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
