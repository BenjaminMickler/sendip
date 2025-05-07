# SendIP

Consists of a simple client that sends its local IP address to a server on boot. The server displays the IP address on a webpage. Useful for SSHing into headless devices on different networks or if the local IP changes regularly.

## Building and installing

The client is the headless device and the server should have a static public IP address.

Clone this repo on the client and the server.
```
git clone https://github.com/BenjaminMickler/sendip.git
```

Install any dependencies.
```
go mod tidy
```

Then build the client or the server depending on the device.
```
go build ./client/main.go -o sendip
```
**OR**
```
go build ./server/main.go -o sendip
```

Copy the executable to `/usr/local/bin`.
```
cp sendip /usr/local/bin
```

Both the client and the server need a config file.

Client:
```toml
server = "YOUR_SERVER_IP"
port = "PORT"
```

Server:
```toml
port = "PORT"
```

SendIP will look for `/boot/sendip.toml`, `/etc/sendip.toml` and `./sendip-[server OR client].toml` (for testing). The name and location of the config file is not configurable.

Install the systemd service (same for client and server) and enable it.
```
sudo cp sendip.service /etc/systemd/system/
sudo systemctl enable --now sendip
```

## HTTPS

A reverse proxy is recommended. The following instructions are for NGINX and Let's Encrypt.

TODO

## TODO

- authentication
- multiple clients (with names)