FROM golang:latest

WORKDIR /wloop/relay_server

RUN \
  go get github.com/lucas-clemente/quic-go && \
	go get github.com/waterloop/wcomms/wjson && \
  go get github.com/waterloop/wpool && \
	go get github.com/waterloop/wstream

COPY . .

ENV \
  POD_PORT=42000 \
  DATA_PORT=42000 \
  COMMAND_PORT=42001

EXPOSE $POD_PORT $DATA_PORT $COMMAND_PORT

CMD ["go", "run", "relay_server.go"]