export default function OpenChatSocket() {
  if (socket !== undefined) {
    return;
  }
  socket = new WebSocket("ws://localhost:8080/api/ws");
  console.log("Attempting Websocket Connetcion");
  socket.onmessage = function (evt) {
    const eventData = JSON.parse(evt.data);
    const event = Object.assign(new Event(eventData));
    routeEvent(event);
  };
  socket.onopen = () => {
    console.log("Connected");
  };

  socket.onclose = (event) => {
    console.log("Socket Closed", event);
    let name = "session_token";
    document.cookie = name + "=; Max-Age=-99999999;";
    console.log("closed and out");
  };

  socket.onerror = (error) => {
    console.log("Socket Error:", error);
  };
}