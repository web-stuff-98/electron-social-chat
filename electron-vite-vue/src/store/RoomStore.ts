import { reactive } from "vue";
import { IRoom, IRoomCard } from "../interfaces/GeneralInterfaces";
import { getRoom, getRoomDisplayData } from "../services/Rooms";
import { socketStore } from "./SocketStore";

/**
 * Interval runs inside App.vue which uses the "disappearedRooms" array to
 * clear the room from the "rooms" cache after the room has been "disappeared"
 * for longer than 30 seconds
 */

type DisappearedRoom = {
  ID: string;
  lastSeen: number;
};

interface IRoomStore {
  rooms: Partial<IRoom>[];
  currentRoom: string;

  visibleRooms: string[];
  disappearedRooms: DisappearedRoom[];

  roomEnteredView: (id: string, full?: boolean) => Promise<Partial<IRoom>>;
  roomLeftView: (id: string) => void;

  getRoom: (id: string) => Partial<IRoom> | undefined;

  cacheRoomData: (
    id: string,
    full?: boolean,
    force?: boolean
  ) => Promise<Partial<IRoom>>;
}

export const roomStore: IRoomStore = reactive({
  rooms: [],
  currentRoom: "",

  visibleRooms: [],
  disappearedRooms: [],

  roomEnteredView: async (id: string, full?: boolean) => {
    roomStore.disappearedRooms = roomStore.disappearedRooms.filter(
      (r) => r.ID !== id
    );
    roomStore.visibleRooms = [...roomStore.visibleRooms, id];
    socketStore.send(
      JSON.stringify({
        event_type: "WATCH_ROOM",
        ID: id,
      })
    );
    return roomStore.cacheRoomData(id, full, full);
  },
  roomLeftView: (id: string) => {
    const i = roomStore.visibleRooms.findIndex((u) => u === id);
    if (i === -1) return;
    roomStore.visibleRooms.splice(i, 1);
    if (roomStore.visibleRooms.findIndex((r) => r === id) === -1) {
      socketStore.send(
        JSON.stringify({ event_type: "STOP_WATCHING_ROOM", ID: id })
      );
    }
  },

  getRoom: (id: string) => roomStore.rooms.find((r) => r.ID === id),

  cacheRoomData: async (id: string, full?: boolean, force?: boolean) => {
    const found = roomStore.rooms.find((u) => u.ID === id);
    if (found && !force) return found;
    try {
      const p = full ? getRoom(id) : getRoomDisplayData(id);
      const r = await p;
      roomStore.rooms = [...roomStore.rooms.filter((r) => r.ID !== id), r];
      return r;
    } catch (e) {
      console.warn(`Failed to cache room display data for ${id}: ${e}`);
    }
  },
});
