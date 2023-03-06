<script setup lang="ts">
import {
  welcomeModalStore,
  EWelcomeModalType,
} from "../../../store/WelcomeModalStore";
</script>

<template>
  <h1>Electron Social Chat</h1>
  <p>
    This is my electron, Vue & Golang social app. You may log in using an
    existing example account or register a new one.
  </p>
  <div class="buttons">
    <button @click="welcomeModalStore.modalType = EWelcomeModalType.LOGIN">
      Login
    </button>
    <button @click="welcomeModalStore.modalType = EWelcomeModalType.REGISTER">
      Register
    </button>
  </div>
</template>

<style lang="scss" scoped>
h1 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
}
p {
  margin: 0;
}
.buttons {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  gap: var(--padding-medium);
  button {
    margin: 0;
    flex-grow: 1;
    padding: 0 var(--padding-medium);
  }
}
</style>
