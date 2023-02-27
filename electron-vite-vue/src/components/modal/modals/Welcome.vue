<script lang="ts">
import { modalStore, EModalType } from "../../../store/ModalStore";
export default {
  data() {
    return {
      modalStore,
      EModalType,
    };
  },
};
</script>

<template>
  <h1>Welcome</h1>
  <p>
    This is my electron, Vue & Golang social app. You may log in using an
    existing example account or register a new one.
  </p>
  <div class="buttons">
    <button @click="modalStore.modalType = EModalType.LOGIN">Login</button>
    <button @click="modalStore.modalType = EModalType.REGISTER">
      Register
    </button>
  </div>
</template>

<style scoped>
h1 {
  margin: 0;
  font-size: 2rem;
}
p {
  margin: 0;
  margin-bottom: var(--padding);
}
.buttons {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  gap: var(--padding-medium);
}
button {
  margin: 0;
  flex-grow: 1;
  padding: 0 var(--padding-medium);
}
</style>
