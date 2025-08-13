import { useEffect, useRef, useState } from 'react';

export default function SocketDemo() {
  const [messages, setMessages] = useState<any[]>([]);
  const [text, setText] = useState('');
  const wsRef = useRef<WebSocket>(null);

  useEffect(() => {
    const ws = new WebSocket('ws://localhost:3000');
    wsRef.current = ws;

    ws.onopen = () => console.log('WS open');
    ws.onmessage = (ev) => {
      try {
        const msg = JSON.parse(ev.data);
        setMessages((m) => [...m, msg]);
      } catch {
        setMessages((m) => [...m, { type: 'raw', payload: ev.data }]);
      }
    };
    ws.onclose = () => console.log('WS closed');

    return () => ws.close(); // cleanup on unmount
  }, []);

  const send = () => {
    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) return;
    wsRef.current.send(text);
    setText('');
  };

  return (
    <div style={{ padding: 16 }}>
      <h2>WebSocket Demo</h2>
      <div style={{ marginBottom: 8 }}>
        <input
          value={text}
          onChange={(e) => setText(e.target.value)}
          placeholder="Type a message…"
        />
        <button onClick={send} disabled={!text}>Send</button>
      </div>
      <pre style={{ background: '#111', color: '#0f0', padding: 12 }}>
        {messages.map((m, i) => `${i + 1}. ${m.type}: ${m.payload}\n`)}
      </pre>
    </div>
  );
}
