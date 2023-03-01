import { reactive } from "vue";
import { IRoomCard } from "../interfaces/GeneralInterfaces";

interface IRoomStore {
  rooms: IRoomCard[];
}

export const roomStore: IRoomStore = reactive({
  rooms: [],
});
