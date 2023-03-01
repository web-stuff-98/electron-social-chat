import { roomStore } from "../store/RoomStore";
const useRoom = (id: string) => roomStore.rooms.find((r) => r.ID === id);
export default useRoom;
