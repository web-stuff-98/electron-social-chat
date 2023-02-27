<script setup lang="ts">
import { EModalType, modalStore } from "../../store/ModalStore";
import Welcome from "./modals/Welcome.vue";
import Login from "./modals/Login.vue";
import Register from "./modals/Register.vue";
import Message from "./modals/Message.vue";
</script>

<template>
  <div class="backdrop" />
  <div class="container">
    <div class="modal">
      <button
        @click="modalStore.modalType = EModalType.WELCOME"
        v-show="
          modalStore.modalType !== EModalType.WELCOME &&
          modalStore.modalType !== EModalType.MESSAGE
        "
        class="close-button"
      >
        <v-icon name="io-close" />
      </button>
      <Welcome v-if="modalStore.modalType === EModalType.WELCOME" />
      <Login v-if="modalStore.modalType === EModalType.LOGIN" />
      <Register v-if="modalStore.modalType === EModalType.REGISTER" />
      <Message v-if="modalStore.modalType === EModalType.MESSAGE" />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.backdrop {
  width: 100%;
  height: 100%;
  position: fixed;
  background: rgba(0, 0, 0, 0.07);
  backdrop-filter: blur(2px);
  filter: opacity(0.9);
}
.container {
  display: flex;
  place-items: center;
  place-content: center;
  width: 100%;
  height: 100%;
  position: fixed;
  z-index: 100;
}
.modal {
  display: flex;
  flex-direction: column;
  gap: var(--padding-medium);
  border: 2px solid var(--base-light);
  background: var(--foreground);
  padding: var(--padding);
  border-radius: var(--border-radius-medium);
  box-shadow: var(--shadow);
  max-width: min(calc(100vw - var(--padding)), 15rem);
  position: relative;
}
.close-button {
  position: absolute;
  background: red;
  top: 2px;
  right: 2px;
  width: 1rem;
  height: 1rem;
  z-index: 99;
  padding: 0;
  border: 2px solid white;
  filter: opacity(0.5);
  display: flex;
  transition: filter 100ms ease;
  svg {
    width: 100%;
    height: 100%;
  }
}
.close-button:hover {
  filter: opacity(1);
}
</style>
