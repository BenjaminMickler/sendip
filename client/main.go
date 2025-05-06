package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	IPAddr string `json:"ipaddr"`
}

var cfg Config

func get_ip() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {
	cfg_file, err := os.ReadFile("/etc/sendip.toml")
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(cfg_file, &cfg)
	if err != nil {
		panic(err)
	}

	ip := get_ip().String()

	_, err = http.Get("http://" + cfg.IPAddr + ":4571/ip?ip=" + ip)
	if err != nil {
		panic(err)
	}
}
