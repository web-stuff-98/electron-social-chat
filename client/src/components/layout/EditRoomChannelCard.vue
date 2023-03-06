<script lang="ts" setup>
import { toRefs, ref } from "vue";
import { IResMsg } from "../../interfaces/GeneralInterfaces";
import useRoomChannel from "../../composables/useRoomChannel";
import MessageModal from "../messageModal/MessageModal.vue";
import { editRoomChannelsData } from "../../store/EditRoomChannelsData";
import { roomStore } from "../../store/RoomStore";

/**
 * Weird complicated looking if checks
 */

const props = defineProps<{
  id: string;
  name?: string;
  mainChannelId: string;
}>();
// if there is no id, it finds the channel from from the insert data using the name instead
const { id, name } = toRefs(props);
const channel = useRoomChannel(id.value, true, name?.value);

const modalConfirmation = ref(() => {});
const modalCancellation = ref<Function | undefined>(() => {});
const showModal = ref(false);
const modalMsg = ref<IResMsg>({ msg: "", err: false, pen: false });
const isEditing = ref(false);

function renderName() {
  if (id.value) {
    // if there is an id value, its a room that already exists
    const found = editRoomChannelsData.update_data.find(
      (c) => c.ID === id.value
    );
    if (found) {
      return found.name;
    } else {
      return channel?.name;
    }
  } else {
    // otherwise its a room not yet created
    const found = editRoomChannelsData.insert_data.find(
      (c) => c.name === name?.value
    );
    return found?.name || "";
  }
}

function getEditValue() {
  const found = id.value
    ? editRoomChannelsData.update_data.find((c) => c.ID === id.value)
    : editRoomChannelsData.insert_data.find((c) => c.name === name?.value);
    if(found) return found.name
    else {
      if(id.value) {
        return channel?.name
      } else {
        return ""
      }
    }
}

function handleEditRoomNameInput(e: Event) {
  const target = e.target as HTMLInputElement;
  if (target.value.length > 24 || !target.value.trim()) return;
  // if there is an id value, its a room that already exists being edited
  if (id.value) {
    const i = editRoomChannelsData.update_data.findIndex(
      (c) => c.ID === id.value
    );
    if (i === -1) {
      editRoomChannelsData.update_data.push({
        ID: id.value,
        name: target.value,
      });
    } else {
      editRoomChannelsData.update_data[i].name = target.value;
    }
  } else {
    // otherwise its a room that hasn't been created yet being edited
    const i = editRoomChannelsData.insert_data.findIndex(
      (c) => c.name === name?.value
    );
    if (i !== -1) {
      editRoomChannelsData.insert_data[i].name = target.value;
    }
  }
}

function handleEditChannelNameForm() {
  isEditing.value = false;
}

function deleteClicked() {
  if (editRoomChannelsData.promote_to_main === id.value && id.value) {
    modalMsg.value = {
      msg: "You cannot delete a channel which you also intend on promoting to main",
      err: true,
      pen: false,
    };
    modalConfirmation.value = () => (showModal.value = false);
    modalCancellation.value = undefined;
    showModal.value = true;
    return;
  }
  const room = roomStore.rooms.find((r) => r.ID === roomStore.currentRoom);
  if (
    room?.main_channel === id.value &&
    !editRoomChannelsData.promote_to_main &&
    !editRoomChannelsData.insert_data.find((c) => c.promote_to_main)
  ) {
    modalMsg.value = {
      msg: "You cannot delete the main channel without replacing it with another",
      err: true,
      pen: false,
    };
    modalConfirmation.value = () => (showModal.value = false);
    modalCancellation.value = undefined;
    showModal.value = true;
    return;
  }
  if (id.value) {
    if (editRoomChannelsData.delete_ids.includes(id.value)) {
      editRoomChannelsData.delete_ids = [
        ...editRoomChannelsData.delete_ids.filter(
          (checkId) => checkId !== id.value
        ),
      ];
    } else {
      editRoomChannelsData.delete_ids = [
        ...editRoomChannelsData.delete_ids.filter(
          (checkId) => checkId !== id.value
        ),
        id.value,
      ];
      modalMsg.value = {
        msg: "Channel is now flagged for deletion",
        err: false,
        pen: false,
      };
      modalConfirmation.value = () => (showModal.value = false);
      modalCancellation.value = undefined;
      showModal.value = true;
    }
  } else {
    // else because if there is no ID, its a channel that hasn't been created,
    // so remove from the insert_data by name instead
    const i = editRoomChannelsData.insert_data.findIndex(
      (c) => c.name === channel?.name
    );
    if (editRoomChannelsData.insert_data[i].promote_to_main) {
      if (editRoomChannelsData.delete_ids.includes(room?.main_channel!)) {
        editRoomChannelsData.delete_ids =
          editRoomChannelsData.delete_ids.filter(
            (c) => c !== room?.main_channel
          );
      }
    }
    editRoomChannelsData.insert_data.splice(i, 1);
  }
}

function promoteToMainClicked() {
  if (editRoomChannelsData.delete_ids.includes(id.value)) {
    modalMsg.value = {
      msg: "You cannot promote a channel to main if is also flagged for deletion",
      err: true,
      pen: false,
    };
    modalConfirmation.value = () => (showModal.value = false);
    modalCancellation.value = undefined;
    showModal.value = true;
    return;
  }
  // if there is an id.value, it means it's a channel that already exists
  if (id.value) {
    if (editRoomChannelsData.promote_to_main === id.value) {
      editRoomChannelsData.promote_to_main = "";
    } else {
      editRoomChannelsData.promote_to_main = id.value;
      modalMsg.value = {
        msg: "Channel will be promoted to main",
        err: true,
        pen: false,
      };
      modalConfirmation.value = () => (showModal.value = false);
      modalCancellation.value = undefined;
      showModal.value = true;
    }
  } else {
    // otherwise it's a channel pending creation, so proceed by name
    editRoomChannelsData.promote_to_main = "";
    const i = editRoomChannelsData.insert_data.findIndex(
      (c) => c.name === channel?.name
    );
    if (!editRoomChannelsData.insert_data[i].promote_to_main) {
      editRoomChannelsData.insert_data = editRoomChannelsData.insert_data.map(
        (c) => ({
          ...c,
          promote_to_main: false,
        })
      );
      editRoomChannelsData.insert_data[i].promote_to_main = true;
    } else {
      editRoomChannelsData.insert_data[i].promote_to_main = false;
      const room = roomStore.rooms.find((r) => r.ID === roomStore.currentRoom);
      if (editRoomChannelsData.delete_ids.includes(room?.main_channel!)) {
        editRoomChannelsData.delete_ids =
          editRoomChannelsData.delete_ids.filter(
            (c) => c !== room?.main_channel
          );
      }
    }
  }
}
</script>

<template>
  <MessageModal
    :msg="modalMsg"
    :show="showModal"
    :confirmationCallback="modalConfirmation"
    :cancellationCallback="modalCancellation"
  />
  <div class="channel">
    <div v-if="!isEditing" class="name">
      {{ renderName() }} {{ mainChannelId === id ? "(main)" : "" }}
    </div>
    <form
      @submit.prevent="handleEditChannelNameForm"
      class="edit-channel-form"
      v-else
    >
      <input
        maxlength="24"
        :value="getEditValue()"
        @input="handleEditRoomNameInput"
        type="text"
      />
      <button type="submit">
        <v-icon name="bi-check-lg" />
      </button>
    </form>
    <div class="buttons">
      <button
        v-if="!isEditing"
        type="button"
        :style="
          editRoomChannelsData.delete_ids.includes(id)
            ? {
                background: 'red',
                border: '1px solid black',
                outline: '1px solid white',
              }
            : {}
        "
        @click="deleteClicked"
      >
        <v-icon name="md-delete-sharp" />
      </button>
      <button v-if="!isEditing" @click="isEditing = true" type="button">
        <v-icon name="ri-edit-box-fill" />
      </button>
      <button
        type="button"
        :style="
          (editRoomChannelsData.promote_to_main === id && id) ||
          editRoomChannelsData.insert_data.find(
            (c) => c.name === channel?.name && c.promote_to_main
          )
            ? { background: 'lime' }
            : {}
        "
        @click="promoteToMainClicked"
        v-if="mainChannelId !== id && !isEditing"
      >
        <v-icon name="bi-caret-up-fill" />
      </button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.channel {
  display: flex;
  justify-content: space-between;
  gap: var(--padding);
  width: 100%;
  box-sizing: border-box;
  align-items: center;
  padding: 2px;
  .edit-channel-form {
    padding: var(--padding-medium);
    gap: 3px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: auto;
    button {
      padding: 0;
      height: 100%;
      display: flex;
      align-items: center;
      margin: 0;
    }
  }
  .name {
    text-align: left;
    padding: 0 3px;
    flex-grow: 1;
    white-space: nowrap;
  }
  .buttons {
    border-radius: var(--border-radius-medium);
    gap: 1px;
    display: flex;
    button {
      border: 1px solid var(--base);
      margin: 0;
      padding: 0;
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: 4px;
      svg {
        width: 1rem;
        height: 1rem;
        padding: 1px;
      }
    }
  }
}
</style>
