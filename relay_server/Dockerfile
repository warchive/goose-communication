FROM golang:latest

WORKDIR /wloop/relay_server

RUN \
  go get -u github.com/lucas-clemente/quic-go && \
	go get -u github.com/waterloop/wcomms/wjson && \
  go get -u github.com/waterloop/wpool && \
	go get -u github.com/waterloop/wstream

COPY . .

ENV \
  POD_PORT=42000 \
  DATA_PORT=42001 \
  COMMAND_PORT=42002

EXPOSE $POD_PORT $DATA_PORT $COMMAND_PORT

CMD ["go", "run", "relay_server.go"]