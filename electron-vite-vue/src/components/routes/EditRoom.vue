<script lang="ts" setup>
import { useRoute } from "vue-router";
import { IRoom, IResMsg } from "../../interfaces/GeneralInterfaces";
import ResMsg from "../layout/ResMsg.vue";
import Toggle from "../shared/Toggle.vue";
import { ref, onMounted, onBeforeUnmount } from "vue";
import {
  getRoom,
  updateRoom,
  updateRoomChannelsData,
} from "../../services/Rooms";
import EditRoomChannelCard from "../layout/EditRoomChannelCard.vue";
import { roomStore } from "../../store/RoomStore";
import { roomChannelStore } from "../../store/RoomChannelStore";
import { editRoomChannelsData } from "../../store/EditRoomChannelsData";

enum EEditRoomSection {
  "BASIC" = "Basic room settings",
  "CHANNELS" = "Room channels",
  "USERS" = "Members & banned users",
}

const section = ref<EEditRoomSection>(EEditRoomSection.BASIC);
const route = useRoute();
const room = ref<IRoom>({
  ID: "",
  name: "",
  blur: "",
  author: "",
  channels: [],
  main_channel: "",
  is_private: false,
});
const originalRoom = ref<IRoom>({
  ID: "",
  name: "",
  blur: "",
  author: "",
  channels: [],
  main_channel: "",
  is_private: false,
});
const resMsg = ref<IResMsg>({ msg: "", err: false, pen: false });

onMounted(async () => {
  editRoomChannelsData.delete_ids = [];
  editRoomChannelsData.insert_data = [];
  editRoomChannelsData.promote_to_main = "";
  editRoomChannelsData.update_data = [];
  try {
    resMsg.value = { msg: "", err: false, pen: true };
    const data: IRoom = await getRoom(route.params.id as string);
    room.value = data;
    originalRoom.value = data;
    await roomChannelStore.getDisplayDataForChannels(route.params.id as string);
    roomStore.rooms = [
      ...roomStore.rooms.filter((r) => r.ID !== data.ID),
      data,
    ];
    roomStore.currentRoom = data.ID;
    resMsg.value = { msg: "", err: false, pen: false };
  } catch (e) {
    resMsg.value = { msg: `${e}`, err: true, pen: false };
  }
});

onBeforeUnmount(async () => {
  roomStore.currentRoom = "";
});

const addChannelInput = ref("");
function handleAddChannelInput(e: Event) {
  const target = e.target as HTMLInputElement;
  if (target.value.length > 24 || !target.value.trim()) return;
  addChannelInput.value = target.value;
}

function handleAddChannelClicked() {
  if (addChannelInput.value.trim() === "" || addChannelInput.value.length > 24)
    return;
  editRoomChannelsData.insert_data.push({
    name: addChannelInput.value,
    promote_to_main: false,
  });
  addChannelInput.value = "";
}

async function handleSubmit() {
  try {
    resMsg.value = { msg: "", err: false, pen: true };
    if (room !== originalRoom)
      await updateRoom(room.value.ID, room.value.name, room.value.is_private);
    if (
      editRoomChannelsData.delete_ids.length !== 0 ||
      editRoomChannelsData.insert_data.length !== 0 ||
      editRoomChannelsData.promote_to_main !== "" ||
      editRoomChannelsData.update_data.length !== 0
    )
      await updateRoomChannelsData(
        room.value.ID,
        editRoomChannelsData.update_data,
        editRoomChannelsData.insert_data,
        editRoomChannelsData.delete_ids,
        editRoomChannelsData.promote_to_main
      );
    resMsg.value = { msg: "", err: false, pen: false };
  } catch (e) {
    resMsg.value = { msg: `${e}`, err: true, pen: false };
  }
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="container">
    <ResMsg :resMsg="resMsg" />
    <div v-if="!resMsg.pen && !resMsg.err && !resMsg.msg" class="content">
      <!-- Basic settings section -->
      <span v-if="section === EEditRoomSection.BASIC">
        <Toggle
          name="Private"
          :on="Boolean(room?.is_private)"
          :toggleFunc="() => (room!.is_private = !room!.is_private)"
        />
        <div class="input-label">
          <label for="name">Room name</label>
          <input v-model="room.name" id="name" type="text" />
        </div>
        <button type="button">Select image</button>
      </span>
      <!-- Channels section -->
      <span v-if="section === EEditRoomSection.CHANNELS">
        <div class="channels-list">
          <div class="channel-container" v-for="channelId in room.channels">
            <EditRoomChannelCard
              :mainChannelId="room.main_channel"
              :id="channelId"
            />
          </div>
          <div
            v-if="editRoomChannelsData.insert_data.length !== 0"
            class="channel-container"
            v-for="channel in editRoomChannelsData.insert_data"
          >
            <EditRoomChannelCard
              :mainChannelId="room.main_channel"
              :id="''"
              :name="channel.name"
            />
          </div>
          <div class="add-channel-container">
            <input
              maxlength="24"
              :value="addChannelInput"
              @input="handleAddChannelInput"
              type="text"
            />
            <v-icon
              @click="handleAddChannelClicked"
              class="add-channel-icon"
              name="io-add-circle-sharp"
            />
          </div>
        </div>
      </span>
      <!-- Users section -->
      <span v-if="section === EEditRoomSection.USERS"> </span>
      <button
        @click="section = EEditRoomSection.BASIC"
        v-if="section !== EEditRoomSection.BASIC"
        type="button"
      >
        Basic settings
      </button>
      <button
        @click="section = EEditRoomSection.CHANNELS"
        v-if="section !== EEditRoomSection.CHANNELS"
        type="button"
      >
        Channels
      </button>
      <button
        @click="section = EEditRoomSection.USERS"
        v-if="section !== EEditRoomSection.USERS"
        type="button"
      >
        Users
      </button>
      <button type="submit">Update room</button>
    </div>
  </form>
</template>

<style lang="scss" scoped>
.container {
  width: 100%;
  height: 100%;
  padding: calc(var(--padding-medium) + 1px);
  box-sizing: border-box;
  display: flex;
  gap: var(--padding-medium);
  flex-direction: column;
  align-items: center;
  justify-content: center;
  span {
    margin: 0;
    gap: var(--padding-medium);
    display: flex;
    flex-direction: column;
  }
  button:first-of-type {
    margin-top: var(--padding-medium);
  }
  input {
    text-align: center;
  }
  .content {
    position: relative;
    box-sizing: border-box;
    width: fit-content;
    height: fit-content;
    border: 2px solid var(--base-light);
    box-shadow: var(--shadow-medium);
    display: flex;
    flex-direction: column;
    border-radius: var(--border-radius-medium);
    padding: var(--padding);
    gap: var(--padding-medium);
    overflow: hidden;
    .channels-list {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: var(--padding-medium);
      box-shadow: var(--shadow-medium);
      border-radius: var(--border-radius-medium);
      border: 1px solid var(--base-light);
      .channel-container {
        width: 100%;
      }
      .add-channel-container {
        display: flex;
        align-items: center;
        gap: 4px;
        padding: 4px;
        padding-right: 0;
        input {
          flex-grow: 1;
          text-align: left;
        }
        .add-channel-icon {
          width: 2rem;
          height: 2rem;
          cursor: pointer;
          padding: 4px;
        }
      }
    }
    .header {
      position: absolute;
      top: 0;
      width: 100%;
      height: 1.75rem;
      display: flex;
      align-items: center;
      justify-content: center;
      text-align: center;
      border-bottom: 1px solid var(--base-light);
      background: rgba(0, 0, 0, 0.166);
    }
  }
}
</style>
