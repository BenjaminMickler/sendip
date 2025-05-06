package main

import (
	"net/http"
)

var ip string

func get_ip(w http.ResponseWriter, r *http.Request) {
	ip = r.URL.Query().Get("ip")
	println("IP Address received: ", ip)
}

func show_ip(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Last received IP address: " + ip))
}

func main() {
	http.HandleFunc("/", show_ip)
	http.HandleFunc("/ip", get_ip)
	http.ListenAndServe(":4571", nil)
}
