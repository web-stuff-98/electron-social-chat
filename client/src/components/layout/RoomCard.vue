<script lang="ts" setup>
import { IResMsg } from "../../interfaces/GeneralInterfaces";
import MessageModal from "../messageModal/MessageModal.vue";
import { toRefs, ref, watch, onMounted, onBeforeUnmount } from "vue";
import { deleteRoom } from "../../services/Rooms";
import { roomStore } from "../../store/RoomStore";
import axios from "axios";

const props = defineProps<{ id: string }>();
const { id } = toRefs(props);

const modalConfirmation = ref(() => {});
const modalCancellation = ref(() => {});
const showModal = ref(false);
const modalMsg = ref<IResMsg>({ msg: "", err: false, pen: false });
const imgObjectURL = ref("");

const containerRef = ref<HTMLCanvasElement | null>(null);

const observer = new IntersectionObserver(async ([entry]) => {
  if (entry.isIntersecting) {
    const r = await roomStore.roomEnteredView(id.value, false);
    const abortController = new AbortController();
    const res = await axios({
      method: "GET",
      url: r.img_url,
      responseType: "arraybuffer",
      withCredentials: true,
    });
    if (imgObjectURL.value) URL.revokeObjectURL(imgObjectURL.value);
    const blob = new Blob([res.data], { type: "image/jpeg" });
    imgObjectURL.value = URL.createObjectURL(blob);
    return () => {
      abortController.abort();
    };
  } else {
    roomStore.roomLeftView(id.value);
  }
});

onMounted(() => {
  observer.observe(containerRef.value!);
});
onBeforeUnmount(() => {
  observer.disconnect();
});

watch(roomStore, async (oldVal, newVal) => {
  const oldRoom = oldVal.rooms.find((r) => r.ID === id.value);
  const newRoom = newVal.rooms.find((r) => r.ID === id.value);
  if (oldRoom?.img_url !== newRoom && !oldRoom?.img_url) {
    const abortController = new AbortController();
    const res = await axios({
      method: "GET",
      url: newRoom?.img_url,
      responseType: "arraybuffer",
      withCredentials: true,
    });
    if (imgObjectURL.value) URL.revokeObjectURL(imgObjectURL.value);
    const blob = new Blob([res.data], { type: "image/jpeg" });
    imgObjectURL.value = URL.createObjectURL(blob);
    return () => {
      abortController.abort();
    };
  }
});

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
  <div
    :style="
      roomStore.getRoom(id)?.blur
        ? {
            'background-image': `url(${
              imgObjectURL || roomStore.getRoom(id)?.blur
            })`,
          }
        : {}
    "
    ref="containerRef"
    class="container"
  >
    <div class="inner">
      <MessageModal
        :msg="modalMsg"
        :show="showModal"
        :confirmationCallback="modalConfirmation"
        :cancellationCallback="modalCancellation"
      />
      <div class="name">
        {{ roomStore.getRoom(id)?.name }}
      </div>
      <div class="buttons">
        <button @click="promptDeleteRoom">
          <v-icon name="md-delete-sharp" />
        </button>
        <router-link :to="`/room/${id}`">
          <button><v-icon name="bi-door-closed-fill" /></button>
        </router-link>
        <router-link :to="`/room/edit/${id}`">
          <button><v-icon name="ri-edit-box-fill" /></button>
        </router-link>
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
  background-position: center;
  text-align: left;
  .inner {
    display: flex;
    justify-content: space-between;
    align-items: center;
    .name {
      color:white;
      font-weight: 600;
      text-shadow: 0px 1px 2px black, 0px 0px 6px black;
    }
    .buttons {
      background: rgba(0, 0, 0, 0.5);
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
