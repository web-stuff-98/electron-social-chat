<script setup lang="ts">
import {
  welcomeModalStore,
  EWelcomeModalType,
} from "../../../store/WelcomeModalStore";
</script>

<template>
  <h1>Tauri-Go-Vue-Chat</h1>
  <span />
  <p>
    This is my Tauri, Vue & Go social app. You may log in using an existing
    example account or register a new one.
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
  font-size: 1.666rem;
  font-weight: 600;
}
span {
  margin: 0;
  padding: 0;
  border-bottom: 2px solid var(--base);
  width: 100%;
  mask-image: linear-gradient(
    to right,
    transparent 0%,
    black 50%,
    transparent 100%
  );
  -webkit-mask-image: linear-gradient(
    to right,
    transparent 0%,
    black 50%,
    transparent 100%
  );
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
    font-size: 1.25rem;
    padding: 0 var(--padding-medium);
  }
}
</style>
