const WebSocket = require('ws');

const ws = new WebSocket('ws://localhost:3000/v1/ws');

ws.on('open', function open() {
    ws.send('Connection Established');
});

ws.on('message', function incoming(data) {
    ws.send(data)
    console.log(data);
});

