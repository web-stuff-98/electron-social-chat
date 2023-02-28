import { reactive } from "vue";

interface ISocketStore {
  socket?: WebSocket;
  connectSocket: (uid: string) => void;
  openSubscription: (name: string) => void;
  openSubscriptions: (names: string[]) => void;
  closeSubscription: (name: string) => void;
  subscriptions: string[];
  send: (data: string) => Boolean;
}

export const socketStore: ISocketStore = reactive({
  socket: undefined,
  connectSocket: (uid: string) => {
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
          event_type: "OPEN_SUBSCRIPTION",
          name: `user=${uid}`,
        })
      );
      socketStore.subscriptions = [`user=${uid}`];
    });
  },
  openSubscription: (name: string) => {
    const sent = socketStore.send(
      JSON.stringify({
        event_type: "OPEN_SUBSCRIPTION",
        name,
      })
    );
    if (sent) {
      socketStore.subscriptions = [
        ...socketStore.subscriptions.filter((sub) => sub !== name),
        name,
      ];
    }
  },
  openSubscriptions: (names) => {
    const sent = socketStore.send(
      JSON.stringify({
        event_type: "OPEN_SUBSCRIPTIONS",
        names,
      })
    );
    if (sent) {
      socketStore.subscriptions = [
        ...socketStore.subscriptions,
        ...names.filter((name) => !socketStore.subscriptions.includes(name)),
      ];
    }
  },
  closeSubscription: (name: string) => {
    const sent = socketStore.send(
      JSON.stringify({
        event_type: "CLOSE_SUBSCRIPTION",
        name,
      })
    );
    if (sent) {
      socketStore.subscriptions = socketStore.subscriptions.filter(
        (sub) => sub !== name
      );
    }
  },
  subscriptions: [],
  send: (data: string): Boolean => {
    if (socketStore.socket && socketStore.socket?.readyState === 1) {
      socketStore.socket.send(data);
      return true;
    }
    console.error("Socket unavailable");
    return false;
  },
});
