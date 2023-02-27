<script setup lang="ts">
import { authStore } from "../../../../store/AuthStore";
import { EModalType, modalStore } from "../../../../store/ModalStore";

import { ref } from "vue";

const fileInputRef = ref<HTMLInputElement>();
defineExpose({ fileInputRef });

function clickPfpHiddenInput() {
  fileInputRef.value?.click();
}

function selectPfp(e: Event) {
  const target = e.target as HTMLInputElement;
  if (!target.files || !target.files[0]) return;
  const file = target.files[0];
  modalStore.modalType = EModalType.MESSAGE;
  modalStore.showModal = true;
  modalStore.messageModalProps = {
    msg: `Are you sure you want to use ${file.name} as your profile picture?`,
    err: false,
    pen: false,
    confirmationCallback: async () => {
      try {
        modalStore.messageModalProps = {
          msg: "Uploading...",
          err: false,
          pen: true,
          confirmationCallback: () => {},
          cancellationCallback: () => {},
        };
        await authStore.uploadPfp(file);
        modalStore.showModal = false;
      } catch (e) {
        modalStore.messageModalProps = {
          msg: `${e}`,
          err: false,
          pen: false,
          confirmationCallback: () => (modalStore.showModal = false),
          cancellationCallback: undefined,
        };
      }
    },
    cancellationCallback: () => {
      modalStore.showModal = false;
    },
  };
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
    <p>Click on your profile picture to select a new one</p>
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
