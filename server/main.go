package main

import (
	"net/http"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Port string `toml:"port"`
}

var cfg Config

var ip string

func get_ip(w http.ResponseWriter, r *http.Request) {
	ip = r.URL.Query().Get("ip")
	println("IP Address received: ", ip)
}

func show_ip(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Last received IP address: " + ip))
}

func main() {
	cfg_path := "/etc/sendip.toml"

	cfg_file, err := os.ReadFile(cfg_path)
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(cfg_file, &cfg)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", show_ip)
	http.HandleFunc("/ip", get_ip)
	http.ListenAndServe(":"+cfg.Port, nil)
}
