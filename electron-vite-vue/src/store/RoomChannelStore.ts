import { reactive } from "vue";
import { IRoomChannel } from "../interfaces/GeneralInterfaces";
import {
  getRoomChannelData,
  getRoomChannelsDisplayData,
} from "../services/Rooms";

interface IRoomChannelStore {
  channels: IRoomChannel[];
  // Get the basic information for an array of channel ids (excludes messages)
  getDisplayDataForChannels: (roomId: string) => Promise<void>;
  // Get the full data for a channel (including messages)
  getFullDataForChannel: (id: string, roomId: string) => Promise<void>;
}

export const roomChannelStore: IRoomChannelStore = reactive({
  channels: [],
  getDisplayDataForChannels: async (roomId: string) => {
    const data: IRoomChannel[] = await getRoomChannelsDisplayData(roomId);
    roomChannelStore.channels = data;
  },
  getFullDataForChannel: async (id: string, roomId: string) => {
    const data: IRoomChannel = await getRoomChannelData(id, roomId);
    roomChannelStore.channels = [
      ...roomChannelStore.channels.filter((r) => r.ID !== id),
      data,
    ];
  },
});
