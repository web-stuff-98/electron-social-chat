<script lang="ts" setup>
import { toRefs, ref } from "vue";
import { IResMsg } from "../../interfaces/GeneralInterfaces";
import useRoomChannel from "../../composables/useRoomChannel";
import MessageModal from "../messageModal/MessageModal.vue";
import { editRoomChannelsData } from "../../store/EditRoomChannelsData";
import { roomStore } from "../../store/RoomStore";
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

function deleteClicked() {
  if (editRoomChannelsData.promoteToMain === id.value && id.value) {
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
    !editRoomChannelsData.promoteToMain &&
    !editRoomChannelsData.insertData.find((c) => c.promoteToMain)
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
    if (editRoomChannelsData.flaggedForDeletion.includes(id.value)) {
      editRoomChannelsData.flaggedForDeletion = [
        ...editRoomChannelsData.flaggedForDeletion.filter(
          (checkId) => checkId !== id.value
        ),
      ];
    } else {
      editRoomChannelsData.flaggedForDeletion = [
        ...editRoomChannelsData.flaggedForDeletion.filter(
          (checkId) => checkId !== id.value
        ),
        id.value,
      ];
      modalMsg.value = {
        msg: "Channel is now flagged for deletion",
        err: true,
        pen: false,
      };
      modalConfirmation.value = () => (showModal.value = false);
      modalCancellation.value = undefined;
      showModal.value = true;
    }
  } else {
    // else because if there is no ID, its a channel that hasn't been created,
    // so remove from the insertData by name instead
    const i = editRoomChannelsData.insertData.findIndex(
      (c) => c.name === channel?.name
    );
    if (editRoomChannelsData.insertData[i].promoteToMain) {
      if (
        editRoomChannelsData.flaggedForDeletion.includes(room?.main_channel!)
      ) {
        editRoomChannelsData.flaggedForDeletion =
          editRoomChannelsData.flaggedForDeletion.filter(
            (c) => c !== room?.main_channel
          );
      }
    }
    editRoomChannelsData.insertData.splice(i, 1);
  }
}

function promoteToMainClicked() {
  if (editRoomChannelsData.flaggedForDeletion.includes(id.value)) {
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
    if (editRoomChannelsData.promoteToMain === id.value) {
      editRoomChannelsData.promoteToMain = "";
    } else {
      editRoomChannelsData.promoteToMain = id.value;
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
    editRoomChannelsData.promoteToMain = "";
    const i = editRoomChannelsData.insertData.findIndex(
      (c) => c.name === channel?.name
    );
    if (!editRoomChannelsData.insertData[i].promoteToMain) {
      editRoomChannelsData.insertData.map((c) => ({
        ...c,
        promoteToMain: false,
      }));
      editRoomChannelsData.insertData[i].promoteToMain = true;
    } else {
      editRoomChannelsData.insertData[i].promoteToMain = false;
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
    <div class="name">
      {{ channel?.name }} {{ mainChannelId === id ? "(main)" : "" }}
    </div>
    <div class="buttons">
      <button
        type="button"
        :style="
          editRoomChannelsData.flaggedForDeletion.includes(id)
            ? { background: 'red' }
            : {}
        "
        @click="deleteClicked"
      >
        <v-icon name="md-delete-sharp" />
      </button>
      <button type="button">
        <v-icon name="ri-edit-box-fill" />
      </button>
      <button
        type="button"
        :style="
          (editRoomChannelsData.promoteToMain === id && id) ||
          editRoomChannelsData.insertData.find(
            (c) => c.name === channel?.name && c.promoteToMain
          )
            ? { background: 'lime' }
            : {}
        "
        @click="promoteToMainClicked"
        v-if="mainChannelId !== id"
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
