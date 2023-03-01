<script lang="ts" setup>
import { authStore } from "../../../../store/AuthStore";
import { ref } from "vue";
import { IResMsg } from "../../../../interfaces/GeneralInterfaces";
import ResMsg from "../../ResMsg.vue";

const fileInputRef = ref<HTMLCanvasElement | null>(null);
const resMsg = ref<IResMsg>({ msg: "", err: false, pen: false });

function clickPfpHiddenInput() {
  fileInputRef.value?.click();
}
async function selectPfp(e: Event) {
  const target = e.target as HTMLInputElement;
  if (!target.files || !target.files[0]) return;
  const file = target.files[0];
  try {
    resMsg.value = { msg: "Uploading...", err: false, pen: true };
    await authStore.uploadPfp(file);
    resMsg.value = { msg: "", err: false, pen: false };
  } catch (e) {
    resMsg.value = { msg: `${e}`, err: true, pen: false };
  }
}
</script>

<template>
  <div class="user">
    <input
      accept=".jpeg,.jpg,.png"
      ref="fileInputRef"
      type="file"
      @change="selectPfp"
    />
    <button
      v-bind:style="{ backgroundImage: `url(${authStore.user?.base64pfp})` }"
      type="button"
      class="pfp"
      @click="clickPfpHiddenInput"
    >
      <v-icon v-if="!authStore.user?.base64pfp" name="la-user" />
    </button>
    <div class="username">
      {{ authStore.user?.username }}
    </div>
    <p v-if="!resMsg.pen">Click on your profile picture to select a new one</p>
    <ResMsg :resMsg="resMsg" />
  </div>
</template>

<style lang="scss" scoped>
.user {
  display: flex;
  flex-direction: column;
  gap: var(--padding-medium);
}

.pfp {
  width: 3rem;
  height: 3rem;
  background: var(--foreground);
  outline: 2px solid var(--base);
  border-radius: 50%;
  box-shadow: var(--shadow-medium);
  cursor: pointer;
  background-size: cover;
  margin: auto;
  svg {
    fill: white;
    width: 80%;
    height: 80%;
  }
}

.username {
  font-weight: 600;
  font-size: 1.1rem;
  text-align: center;
  margin-top: var(--padding-medium);
}

p {
  margin: 0;
  text-align: center;
  font-size: 0.866rem;
}
</style>
