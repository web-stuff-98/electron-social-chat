import { reactive } from "vue";
import { IRoom, IRoomCard } from "../interfaces/GeneralInterfaces";

interface IRoomStore {
  rooms: Partial<IRoom>[];
  currentRoom: string;
}

export const roomStore: IRoomStore = reactive({
  rooms: [],
  currentRoom: "",
});
