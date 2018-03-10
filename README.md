# Backend/Communication-Systems

### Golang installation commands for Raspbian
- Current version: 1.10 (ARM)
- ```wget https://dl.google.com/go/go1.10.linux-armv6l.tar.gz```
- ```sudo tar -C /usr/local -xzf go1.10.linux-armv6l.tar.gz```
- ```export PATH=$PATH:/usr/local/go/bin```

### Golang Dependencies
- ```go get github.com/buger/jsonparser```
- ```go get github.com/lucas-clemente/quic-go```
- ```go get github.com/xtaci/kcp-go```
- ```go get github.com/mogball/wcomms/wjson```

### Backend/Communication-Systems
#### Raspbian Stretch with Desktop Installation Guide
- Version Installed
- Kernel 4.9
- IP Address
- ```10.173.212.248```
- User name
-```pi```
- Password
- ```waterloop```
- Repositories
-```git clone http://prose.io/#teamwaterloop/communication-system-f17/```

#### NodeJS installation commands Raspbian
-Current version: 8.9.0
-``` curl -sL https://deb.nodesource.com/setup_8.x | sudo -E bash ```
- ```sudo apt-get nodejs```

#### Test Network Protocol Speeds to Pi from Laptop over WiFi
- TCP: 18ms-19ms per 100 packets
- UDP: 18ms-19ms per 100 packets
- QUIC (50 Streams Open): ~1.3ms per 100 packets
- Next  Previous

