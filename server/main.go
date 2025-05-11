package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

var cfg Config
var cfg_paths = []string{"/etc/sendip.toml", "./sendip-server.toml"}

type IP struct {
	IP       string
	Name     string
	Time     string
	SendTime time.Time
}

var ips map[string][]IP

var exit chan bool

const static = `
<!DOCTYPE html>
<script>
function copy(text) {
	navigator.clipboard.writeText(text).then(function() {
		console.log('Copied IP address');
	}, function(err) {
		console.error('Could not copy IP address: ', err);
	});
}
</script>

<a href="/">Config generator</a>
<button onclick="location.reload()">Refresh</button>
`

const entry_html = `
<h2>%s</h2>
[%s]: %s <button onclick="copy('%s')">Copy</button>
<hr>
`

const config_template = `host = "%s"
port = "%s"
name = "${name_v}"
key = "${key_v}"`

var config_gen string

func send_ip(w http.ResponseWriter, r *http.Request) {
	var err error

	key, err := url.QueryUnescape(r.URL.Query().Get("key"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name, err := url.QueryUnescape(r.URL.Query().Get("name"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if name == "" {
		name = strings.Split(r.RemoteAddr, ":")[0]
	}

	ip, err := url.QueryUnescape(r.URL.Query().Get("ip"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	time_str, err := url.QueryUnescape(r.URL.Query().Get("time"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key = fmt.Sprintf("%.*s", 100, key)
	ip = fmt.Sprintf("%.*s", 100, ip)
	name = fmt.Sprintf("%.*s", 100, name)
	time_str = fmt.Sprintf("%.*s", 100, time_str)

	ip_struct := IP{
		IP:       ip,
		Name:     name,
		Time:     time_str,
		SendTime: time.Now(),
	}

	if _, ok := ips[key]; ok {
		for i, _ := range ips[key] {
			if ips[key][i].Name == name {
				ips[key][i] = ip_struct
				goto found
			}
		}
		ips[key] = append(ips[key], ip_struct)
	found:
	} else {
		ips[key] = []IP{ip_struct}
	}
	w.WriteHeader(http.StatusOK)
}

func show_ip(w http.ResponseWriter, r *http.Request) {
	key, err := url.QueryUnescape(r.URL.Query().Get("key"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if key == "" {
		w.Write([]byte(config_gen))
		return
	}
	w.Write([]byte(static))
	entries, ok := ips[key]
	if !ok {
		w.Write([]byte("<br>Invalid key"))
		return
	}
	for _, entry := range entries {
		fmt.Fprintf(w, entry_html, entry.Name, entry.Time, entry.IP, entry.IP)
	}
}

func cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	exit = make(chan bool)
	for {
		select {
		case <-exit:
			return
		case <-ticker.C:
			for i, key := range ips {
				for j, ip := range key {
					if time.Since(ip.SendTime) > 12*time.Hour {
						ips[i] = slices.Delete(key, j, j+1)
					}
				}
				if len(key) == 0 {
					delete(ips, i)
				}
			}
		}
	}
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

	config_gen = `
	<!DOCTYPE html>

	<h2>SendIP config generator</h2>
	<label for="name">Name: </label><input id="name" /><br><br>
	<button id="gen_config">Download config file</button>

	<a id="link"></a>

	<p>Tips:
	<ul>
		<li>Copy the config file to the /boot partition or the /etc folder on your device.</li>
		<li>Ensure that the file is named sendip.toml</li>
		<li>The page will show "invalid key" until the SendIP client sends its IP address.</li>
		<li>Bookmark the page in your bookmarks toolbar for quick access.</li>
		<li>All IP addresses are removed ~12 hours after they were received by the server, after this, the page will show "invalid key" until the next IP address is sent.</li>
		<li>Multiple devices can share a key! Simply copy the config file, open it up in a text editor and change the name.</li>
	</ul>
	</p>

	<h2>Get link for key</h2>
	<label for="key">Key: </label><input id="key" /><br><br>
	<a id="key_link"></a>

	<script>
	function download(filename, text) {
		var element = document.createElement("a");
		element.setAttribute("href", "data:text/plain;charset=utf-8," + encodeURIComponent(text));
		element.setAttribute("download", filename);
		element.style.display = "none";
		document.body.appendChild(element);
		element.click();
		document.body.removeChild(element);
	}

	function uuidv4() {
		return "10000000-1000-4000-8000-100000000000".replace(/[018]/g, c =>
			(+c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> +c / 4).toString(16)
		);
	}

	document.getElementById("gen_config").addEventListener("click", function() {
		var name_v = document.getElementById("name").value;
		var key_v = uuidv4();
		document.getElementById("link").href = ` + "`${location.protocol}//" + cfg.Host + ":" + cfg.Port + "/?key=${key_v}`" + `;
		document.getElementById("link").innerHTML = "Use this link to view your local IP address";
		download("sendip.toml", ` + "`" + fmt.Sprintf(config_template, cfg.Host, cfg.Port) + "`" + `);
	});

	document.getElementById("key").addEventListener("input", function() {
		var link = document.getElementById("key_link");
		var key_v = document.getElementById("key").value;
		link.href = ` + "`${location.protocol}//" + cfg.Host + ":" + cfg.Port + "/?key=${key_v}`" + `;
		link.innerHTML = ` + "`${location.protocol}//" + cfg.Host + ":" + cfg.Port + "/?key=${key_v}`" + `;
	});
	</script>
	`

	ips = make(map[string][]IP)

	go cleanup()

	http.HandleFunc("/", show_ip)
	http.HandleFunc("/sendip", send_ip)
	http.ListenAndServe(":"+cfg.Port, nil)
}
