# SendIP

Consists of a simple client that sends its local IP address to a server on boot. The server displays the IP address on a webpage. Useful for SSHing into headless devices on different networks or if the local IP changes regularly.

## Building and installing

Clone this repo on the client and the server. The client is the headless device and the server should have a static public IP address.
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
# OR
go build ./server/main.go -o sendip
```

Copy the executable to `/usr/local/bin`.
```
cp sendip /usr/local/bin
```

Create a config file (**client only**) containing the IP address of the server.
```
sudo echo "ipaddr = 'YOUR_SERVER_IP_ADDRESS'" > /etc/sendip.toml
```

Install the systemd service (same for client and server) and enable it.
```
sudo cp sendip.service /etc/systemd/system/
sudo systemctl enable --now sendip
```

## TODO

- authentication
- multiple clients (with names)
- timestamps