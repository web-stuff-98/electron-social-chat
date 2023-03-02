<script lang="ts" setup>
import { IResMsg } from "../../interfaces/GeneralInterfaces";
import MessageModal from "../messageModal/MessageModal.vue";
import { toRefs, ref } from "vue";
import useRoomCard from "../../composables/useRoomCard";
import { deleteRoom } from "../../services/Rooms";

const props = defineProps<{ id: string }>();
const { id } = toRefs(props);
const room = useRoomCard(id.value);

const modalConfirmation = ref(() => {});
const modalCancellation = ref(() => {});
const showModal = ref(false);
const modalMsg = ref<IResMsg>({ msg: "", err: false, pen: false });

function promptDeleteRoom() {
  modalMsg.value = {
    msg: "Are you sure you want to delete this room?",
    err: false,
    pen: false,
  };
  showModal.value = true;
  modalConfirmation.value = async () => {
    try {
      modalMsg.value = {
        msg: "Deleting room...",
        err: false,
        pen: true,
      };
      await deleteRoom(id.value);
      showModal.value = false;
    } catch (e) {
      modalMsg.value = {
        msg: `${e}`,
        err: true,
        pen: false,
      };
    }
  };
  modalCancellation.value = () => (showModal.value = false);
}
</script>

<template>
  <MessageModal
    :msg="modalMsg"
    :show="showModal"
    :confirmationCallback="modalConfirmation"
    :cancellationCallback="modalCancellation"
  />
  <div
    :style="room?.blur ? { 'background-image': `url(${room.blur})` } : {}"
    class="container"
  >
    <div class="inner">
      <div class="name">
        {{ room?.name }}
      </div>
      <div class="buttons">
        <button @click="promptDeleteRoom">
          <v-icon name="md-delete-sharp" />
        </button>
        <router-link :to="`/room/${room!.ID}`">
          <button><v-icon name="bi-door-closed-fill" /></button>
        </router-link>
        <button><v-icon name="ri-edit-box-fill" /></button>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.container {
  width: 100%;
  min-width: 100%;
  padding: 1px;
  padding-left: 4px;
  box-sizing: border-box;
  background-size: cover;
  text-align: left;
  .inner {
    display: flex;
    justify-content: space-between;
    align-items: center;
    .name {
      text-shadow: 0px 1px 2px black, 0px 0px 6px black;
    }
    .buttons {
      background: rgba(0, 0, 0, 0.5);
      backdrop-filter: blur(2px);
      border: 1px solid var(--base);
      display: flex;
      align-items: center;
      padding: 1px;
      border-radius: var(--border-radius-medium);
      gap: 1px;
      button {
        border: 1px solid var(--base);
        padding: 0;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 1px;
        border-radius: 4px;
        svg {
          width: 1rem;
          height: 1rem;
        }
      }
    }
  }
}
</style>
