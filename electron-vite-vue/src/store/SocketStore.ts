import { reactive } from "vue";

interface ISocketStore {
  socket?: WebSocket;
  openSubscription: (name: string) => void;
  openSubscriptions: (names: string[]) => void;
  closeSubscription: (name: string) => void;
  subscriptions: string[];
  send: (data: string) => Boolean;
}

export const socketStore: ISocketStore = reactive({
  socket: undefined,
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
