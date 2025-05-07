package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Port string `toml:"port"`
}

var cfg Config
var cfg_paths = []string{"/etc/sendip.toml", "./sendip-server.toml"}

var ip string
var time_str string

const static = `
<!DOCTYPE html>
[%s]: %s <button id="copy">Copy</button>
<script>
document.getElementById("copy").addEventListener("click", function() {
	navigator.clipboard.writeText("%s").then(function() {
		console.log('Copied IP address');
	}, function(err) {
		console.error('Could not copy IP address: ', err);
	});
});
</script>
`

func get_ip(w http.ResponseWriter, r *http.Request) {
	var err error
	ip, err = url.QueryUnescape(r.URL.Query().Get("ip"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	time_str, err = url.QueryUnescape(r.URL.Query().Get("time"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	println("[" + time_str + "]: " + ip)
	w.WriteHeader(http.StatusOK)
}

func show_ip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, static, time_str, ip, ip)
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

	http.HandleFunc("/", show_ip)
	http.HandleFunc("/sendip", get_ip)
	http.ListenAndServe(":"+cfg.Port, nil)
}
