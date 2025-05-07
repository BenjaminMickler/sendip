package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	IPAddr string `toml:"ipaddr"`
	Port   string `toml:"port"`
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
	cfg_path := "/boot/sendip.toml"
	if _, err := os.Stat(cfg_path); errors.Is(err, os.ErrNotExist) {
		cfg_path = "/etc/sendip.toml"
	}

	cfg_file, err := os.ReadFile(cfg_path)
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(cfg_file, &cfg)
	if err != nil {
		panic(err)
	}

	ip := get_ip().String()

	_, err = http.Get("http://" + cfg.IPAddr + ":" + cfg.Port + "/ip?ip=" + ip)
	if err != nil {
		panic(err)
	}
}
