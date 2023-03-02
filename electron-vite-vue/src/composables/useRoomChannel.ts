import { roomChannelStore } from "../store/RoomChannelStore";
const useRoomChannel = (id: string) =>
  roomChannelStore.channels.find((c) => c.ID === id);
export default useRoomChannel;
