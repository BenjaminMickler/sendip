package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Server string `toml:"server"`
	Port   string `toml:"port"`
}

var cfg Config
var cfg_paths = []string{"/boot/sendip.toml", "/etc/sendip.toml", "./sendip-client.toml"}

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
	cfg_i := 0

next_cfg:
	cfg_path := cfg_paths[cfg_i]
	if _, err := os.Stat(cfg_path); errors.Is(err, os.ErrNotExist) {
		cfg_i += 1
		if cfg_i >= len(cfg_paths) {
			panic(errors.New("no config files found"))
		}
		goto next_cfg
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
	time_str := time.Now().Format("2 Jan 2006 15:04:05")
	req_url := "http://" + cfg.Server + ":" + cfg.Port + "/ip?ip=" + url.QueryEscape(ip) + "&time=" + url.QueryEscape(time_str)
	_, err = http.Get(req_url)
	if err != nil {
		panic(err)
	}
}
