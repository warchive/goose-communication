# Websockets

Pull client dependencies
```
make setup-client
```

Build the server
```
make build-server
```

Run the server
```
make run-server
```

Run the client
```
make run-client
```

This will start bi-directional communication with server and the client
over websockets. The server sends a packet every 1 second and the client
replies to that message by sending another packet. Both packets are logged
to console.