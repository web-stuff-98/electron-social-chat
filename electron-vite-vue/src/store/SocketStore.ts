import { reactive } from "vue";

interface ISocketStore {
    socket?: WebSocket
}

export const socketStore:ISocketStore = reactive({
    socket: undefined,
})