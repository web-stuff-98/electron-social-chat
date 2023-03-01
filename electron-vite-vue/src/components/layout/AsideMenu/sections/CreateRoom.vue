<script lang="ts" setup>
import Toggle from "../../../shared/Toggle.vue";
import ResMsg from "../../ResMsg.vue";
import { createRoom, uploadRoomImage } from "../../../../services/Rooms";
import { ref } from "vue";
import { IResMsg } from "../../../../interfaces/GeneralInterfaces";

const fileInputRef = ref<HTMLCanvasElement | null>(null);

const roomName = ref("");
const roomPrivate = ref(false);
const roomImage = ref<File>();
const roomImageURL = ref<string>();

const resMsg = ref<IResMsg>({ msg: "", err: false, pen: false });

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

async function handleSubmit() {
  try {
    resMsg.value = { msg: "", err: false, pen: true };
    const id = await createRoom(roomName.value, roomPrivate.value);
    if (roomImage.value) {
      await uploadRoomImage(roomImage.value, id);
    }
    resMsg.value = { msg: "", err: false, pen: false };
  } catch (e) {
    resMsg.value = { msg: `${e}`, err: true, pen: false };
  }
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="container">
    <Toggle
      name="Private"
      :on="roomPrivate"
      :toggleFunc="() => (roomPrivate = !roomPrivate)"
    />
    <div class="input-label">
      <label for="name">Room name</label>
      <input v-model="roomName" id="name" type="text" />
    </div>
    <input
      accept=".jpeg,.jpg,.png"
      ref="fileInputRef"
      type="file"
      @change="selectImage"
    />
    <button @click="clickHiddenImageInput" type="button">Select image</button>
    <button type="submit">Create room</button>
    <img v-if="roomImageURL" :src="roomImageURL" />
    <ResMsg :resMsg="resMsg" />
  </form>
</template>

<style lang="scss" scoped>
.container {
  max-width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--padding-medium);
  padding: var(--padding-medium);
  input {
    width: 100%;
    text-align: center;
  }
  button {
    width: 100%;
  }
  img {
    width: 100%;
    object-position: center;
    object-fit: contain;
    max-height: 5rem;
  }
}
</style>
