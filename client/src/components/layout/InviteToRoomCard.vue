<script lang="ts" setup>
import axios from "axios";
import { onMounted, onBeforeUnmount, ref, toRefs, watch } from "vue";
import { roomStore } from "../../store/RoomStore";
const props = defineProps<{ id: string }>();

const { id } = toRefs(props);
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

onMounted(() => {
  observer.observe(containerRef.value!);
});
onBeforeUnmount(() => {
  observer.disconnect();
});
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
    class="card"
  >
    {{ roomStore.getRoom(id)?.name }}
  </div>
</template>

<style lang="scss" scoped>
.card {
  display: flex;
  padding: var(--padding-medium);
  box-sizing: border-box;
  background-size: cover;
  background-position: center;
  text-align: left;
  text-shadow: var(--shadow-medium);
  font-weight: 600;
  cursor: pointer;
}
</style>
