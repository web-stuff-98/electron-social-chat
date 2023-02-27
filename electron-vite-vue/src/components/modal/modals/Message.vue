<script setup lang="ts">
import { modalStore } from "../../../store/ModalStore";
</script>

<template>
  <v-icon
    v-if="modalStore.messageModalProps.pen"
    name="pr-spinner"
    class="spinner anim-spin"
  />
  <v-icon
    v-if="modalStore.messageModalProps.err"
    class="err-icon"
    name="md-error-round"
  />
  <p>
    {{ modalStore.messageModalProps.msg }}
  </p>
  <div class="buttons">
    <button
      v-if="
        modalStore.messageModalProps.cancellationCallback &&
        !modalStore.messageModalProps.pen
      "
      @click="modalStore.messageModalProps.cancellationCallback"
    >
      Cancel
    </button>
    <button
      v-if="
        modalStore.messageModalProps.confirmationCallback &&
        !modalStore.messageModalProps.pen
      "
      @click="modalStore.messageModalProps.confirmationCallback"
    >
      Confirm
    </button>
  </div>
</template>

<style scoped>
p {
  margin: 0;
}
.buttons {
  display: flex;
  gap: var(--padding-medium);
  width: 100%;
}
.buttons button {
  flex-grow: 1;
}
.spinner {
  width: 2.5rem;
  height: 2.5rem;
  margin: auto;
}
.err-icon {
  width: 1.5rem;
  height: 1.5rem;
  margin: auto;
}
</style>
