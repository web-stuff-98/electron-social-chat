import { reactive } from "vue";

interface ISocketStore {
  socket?: WebSocket;
  connectSocket: (uid: string) => void;
  send: (data: string) => Boolean;
}

export const socketStore: ISocketStore = reactive({
  socket: undefined,
  connectSocket: (uid: string) => {
    if (socketStore.socket) socketStore.socket.close();
    const socket = new WebSocket(
      process.env.NODE_ENV === "development" ||
      window.location.origin === "http://localhost:8080"
        ? "ws://localhost:8080/api/ws"
        : "wss://electron-social-chat-backend.herokuapp.com/api/ws"
    );
    socketStore.socket = socket;
    socket.addEventListener("open", () => {
      socket.send(
        JSON.stringify({
          event_type: "WATCH_USER",
          ID: uid,
        })
      );
    });
  },
  send: (data: string): Boolean => {
    if (socketStore.socket && socketStore.socket?.readyState === 1) {
      socketStore.socket.send(data);
      return true;
    }
    console.error("Socket unavailable");
    return false;
  },
});
