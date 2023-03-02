import { roomStore } from "../store/RoomStore";
const useRoomCard = (id: string) => roomStore.rooms.find((r) => r.ID === id);
export default useRoomCard;
