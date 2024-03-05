const ws = new WebSocket("ws://localhost:8000/ws");

ws.onmessage = (msg) => {
    console.log(msg);
}