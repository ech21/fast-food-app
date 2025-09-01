import { WebSocketServer } from 'ws';

const wss = new WebSocketServer({ port: 3001 });
console.log('WS server on ws://localhost:3001');

wss.on('connection', (socket) => {
  console.log('Client connected'); 
  socket.send(JSON.stringify({ type: 'hello', payload: 'Welcome!' }));

  socket.on('message', (data) => {
    // echo + broadcast
    for (const client of wss.clients) {
      if (client.readyState === 1) {
        client.send(JSON.stringify({ type: 'msg', payload: data.toString() }));
      }
    }
  });

  socket.on('close', () => console.log('Client disconnected'));
});


//to be: the server side websocket
