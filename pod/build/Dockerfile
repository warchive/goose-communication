FROM golang:latest

WORKDIR /wloop/pod_server

RUN \
  go get -u github.com/lucas-clemente/quic-go && \
	go get -u github.com/waterloop/wcomms/wjson && \
	go get -u github.com/waterloop/wstream

COPY . .

ENV \
  POD_PORT=42000

EXPOSE $POD_PORT

CMD ["go", "run", "pod_server.go"]