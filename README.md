# Backend/Communication-Systems

### Golang installation commands for Raspbian
- Current version: 1.9.2 (ARM)
- ```wget https://storage.googleapis.com/golang/go1.9.2.linux-armv6l.tar.gz```
- ```sudo tar -C /usr/local -xzf go1.9.2.linux-armv6l.tar.gz```
- ```export PATH=$PATH:/usr/local/go/bin```

### Golang Dependencies
- ```go get github.com/buger/jsonparser```
- ```go get github.com/lucas-clemente/quic-go```
- ```go get github.com/xtaci/kcp-go```
