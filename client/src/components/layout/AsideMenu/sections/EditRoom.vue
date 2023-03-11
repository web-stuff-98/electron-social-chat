<script lang="ts" setup>
import { IRoom, IResMsg } from "../../../../interfaces/GeneralInterfaces";
import ResMsg from "../../ResMsg.vue";
import Toggle from "../../../shared/Toggle.vue";
import { onMounted, ref } from "vue";
import {
  updateRoom,
  updateRoomChannelsData,
  uploadRoomImage,
} from "../../../../services/Rooms";
import EditRoomChannelCard from "../../EditRoomChannelCard.vue";
import { roomChannelStore } from "../../../../store/RoomChannelStore";
import { editRoomChannelsData } from "../../../../store/EditRoomChannelsData";
import { roomStore } from "../../../../store/RoomStore";
import { editRoomId } from "../../../../store/CreateEditRoomStore";

enum EEditRoomSection {
  "BASIC" = "Basic room settings",
  "CHANNELS" = "Room channels",
  "USERS" = "Members & banned users",
}

const section = ref<EEditRoomSection>(EEditRoomSection.BASIC);
const room = ref<IRoom>({
  ID: "",
  name: "",
  blur: "",
  author: "",
  channels: [],
  main_channel: "",
  is_private: false,
});
const resMsg = ref<IResMsg>({ msg: "", err: false, pen: false });

const roomImage = ref<File>();
const roomImageURL = ref<string>();
const fileInputRef = ref<HTMLCanvasElement | null>(null);
function clickHiddenImageInput() {
  fileInputRef.value?.click();
}
function selectImage(e: Event) {
  const target = e.target as HTMLInputElement;
  if (!target.files || !target.files[0]) return;
  const file = target.files[0];
  roomImage.value = file;
  roomImageURL.value = URL.createObjectURL(file);
}

const addChannelInput = ref("");
function handleAddChannelInput(e: Event) {
  const target = e.target as HTMLInputElement;
  if (target.value.length > 24 || !target.value.trim()) return;
  addChannelInput.value = target.value;
}

onMounted(async () => {
  const abortController = new AbortController();
  editRoomChannelsData.delete_ids = [];
  editRoomChannelsData.insert_data = [];
  editRoomChannelsData.promote_to_main = "";
  editRoomChannelsData.update_data = [];
  try {
    resMsg.value = { msg: "", err: false, pen: true };
    const data = await roomStore.cacheRoomData(editRoomId.value, true, true);
    room.value = data as IRoom;
    await roomChannelStore.getDisplayDataForChannels(editRoomId.value);
    resMsg.value = { msg: "", err: false, pen: false };
  } catch (e) {
    resMsg.value = { msg: `${e}`, err: true, pen: false };
  }
  return () => {
    abortController.abort();
  };
});

function handleAddChannelClicked() {
  if (addChannelInput.value.trim() === "" || addChannelInput.value.length > 24)
    return;
  editRoomChannelsData.insert_data.push({
    name: addChannelInput.value,
    promote_to_main: false,
  });
  addChannelInput.value = "";
}

function handleRoomNameInput(e: Event) {
  const target = e.target as HTMLInputElement;
  if (target.value.length > 24 || !target.value.trim()) {
    return;
  }
  room.value.name = target.value;
}

async function handleSubmit() {
  const abortController = new AbortController();
  try {
    resMsg.value = { msg: "", err: false, pen: true };
    await updateRoom(room.value.ID, room.value.name, room.value.is_private);
    if (
      editRoomChannelsData.delete_ids.length !== 0 ||
      editRoomChannelsData.insert_data.length !== 0 ||
      editRoomChannelsData.promote_to_main !== "" ||
      editRoomChannelsData.update_data.length !== 0
    ) {
      const insertedRooms = await updateRoomChannelsData(
        room.value.ID,
        editRoomChannelsData.update_data,
        editRoomChannelsData.insert_data,
        editRoomChannelsData.delete_ids,
        editRoomChannelsData.promote_to_main
      );
      roomChannelStore.channels.filter(
        (c) => !editRoomChannelsData.delete_ids.includes(c.ID)
      );
      editRoomChannelsData.update_data.forEach((c) => {
        const i = roomChannelStore.channels.findIndex((ec) => ec.ID === c.ID);
        roomChannelStore.channels[i].name = c.name;
      });
      roomChannelStore.channels =
        roomChannelStore.channels.concat(insertedRooms);
      if (editRoomChannelsData.promote_to_main) {
        room.value.main_channel = editRoomChannelsData.promote_to_main;
      }
      editRoomChannelsData.delete_ids = [];
      editRoomChannelsData.insert_data = [];
      editRoomChannelsData.update_data = [];
      editRoomChannelsData.promote_to_main = "";
    }
    if (roomImage.value) {
      await uploadRoomImage(roomImage.value, editRoomId.value);
    }
    resMsg.value = { msg: "", err: false, pen: false };
  } catch (e) {
    resMsg.value = { msg: `${e}`, err: true, pen: false };
  }
  return () => {
    abortController.abort();
  };
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="container">
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
          <input
            accept=".jpeg,.jpg,.png"
            ref="fileInputRef"
            type="file"
            @change="selectImage"
          />
          <input
            maxlength="24"
            :value="room.name"
            @input="handleRoomNameInput"
            id="name"
            type="text"
          />
        </div>
        <button @click="clickHiddenImageInput" type="button">
          Select image
        </button>
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
      <span class="users" v-if="section === EEditRoomSection.USERS">
        <div class="list-container">
          <label>Banned</label>
          <div class="users-banned"></div>
        </div>
        <div class="list-container">
          <label>Members</label>
          <div class="users-members"></div>
        </div>
      </span>
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
      <img v-if="roomImageURL" :src="roomImageURL" />
    </div>
    <ResMsg :resMsg="resMsg" />
  </form>
</template>

<style lang="scss" scoped>
.container {
  width: 100%;
  height: 100%;
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
    box-shadow: var(--shadow-medium);
    border-radius: var(--border-radius-medium);
    border: 1px solid var(--base-light);
    padding: var(--padding-medium);
  }
  button:first-of-type {
    margin-top: var(--padding-medium);
  }
  input {
    text-align: center;
  }
  img {
    width: 100%;
    object-position: center;
    object-fit: contain;
    max-height: 5rem;
    margin-top: var(--padding-medium);
  }
  .content {
    position: relative;
    box-sizing: border-box;
    width: fit-content;
    height: fit-content;
    display: flex;
    flex-direction: column;
    border-radius: var(--border-radius-medium);
    padding: var(--padding);
    gap: var(--padding-medium);
    overflow: hidden;
    .users {
      display: flex;
      gap: var(--padding-medium);
      flex-direction: row;
      .list-container {
        display: flex;
        flex-direction: column;
        label {
          margin: 0;
          padding: 0;
        }
        .users-banned,
        .users-members {
          min-width: 5rem;
          border: 1px solid var(--base-light);
          min-height: 5rem;
          border-radius: var(--border-radius-medium);
        }
      }
    }
    .channels-list {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: var(--padding-medium);
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
      background: rgba(0, 0, 0, 0.1);
    }
  }
}
</style>
