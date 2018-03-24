FROM golang:latest

WORKDIR /wloop/pod_server

RUN \
  go get github.com/lucas-clemente/quic-go && \
	go get github.com/waterloop/wcomms/wjson && \
	go get github.com/waterloop/wstream

COPY . .

ENV \
  POD_PORT=42000

EXPOSE $POD_PORT

CMD ["go", "run", "pod_server.go"]