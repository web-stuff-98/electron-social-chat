import { roomChannelStore } from "../store/RoomChannelStore";
import { editRoomChannelsData } from "../store/EditRoomChannelsData";

/**
 * If using isEditPage, name must be defined
 */

const useRoomChannel = (
  identifier: string,
  isEditPage?: boolean,
  name?: string
) => {
  const found = roomChannelStore.channels.find((c) => c.ID === identifier);
  if (found && identifier) return found;
  if (!found && isEditPage) {
    return editRoomChannelsData.insert_data.find((c) => c.name === name);
  }
};
export default useRoomChannel;
