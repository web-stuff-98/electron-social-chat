<script setup lang="ts">
import { ipcRenderer } from "electron";
import { ref } from "vue";
import { authStore } from "../../store/AuthStore";
import { modalStore, EModalType } from "../../store/ModalStore";

import { onBeforeMount, onBeforeUnmount } from "vue";

const showAccountMenu = ref(false);
const mouseInAccountMenu = ref(false);

function quit() {
  ipcRenderer.send("window", ["QUIT"]);
}
function minimize() {
  ipcRenderer.send("window", ["MINIMIZE"]);
}
function click() {
  if (!mouseInAccountMenu.value) showAccountMenu.value = false;
}
function deleteAccount() {
  modalStore.modalType = EModalType.MESSAGE;
  modalStore.messageModalProps = {
    msg: "Are you sure you want to delete your account?",
    err: false,
    pen: false,
    confirmationCallback: () => authStore.deleteAccount(),
    cancellationCallback: () => {
      modalStore.showModal = false;
    },
  };
  modalStore.showModal = true;
}

onBeforeMount(() => {
  window.document.addEventListener("click", click);
});
onBeforeUnmount(() => {
  window.document.removeEventListener("click", click);
});
</script>

<template>
  <header>
    <div v-if="authStore.user" class="user-name">
      {{ authStore.user?.username }}
    </div>
    <div class="buttons">
      <div
        @mouseenter="mouseInAccountMenu = true"
        @mouseleave="mouseInAccountMenu = false"
        v-show="showAccountMenu && !modalStore.showModal"
        class="account-menu"
      >
        <button @click="deleteAccount">Delete account</button>
        <button @click="authStore.logout">Log out</button>
      </div>
      <button
        v-if="authStore.user"
        @mouseenter="mouseInAccountMenu = true"
        @mouseleave="mouseInAccountMenu = false"
        @click="showAccountMenu = !showAccountMenu"
        class="account-button"
      >
        <v-icon name="md-manageaccounts-sharp" />
      </button>
      <button @click="minimize" class="minimize-button">
        <v-icon name="fa-minus" />
      </button>
      <button @click="quit" class="quit-button">
        <v-icon name="io-close" />
      </button>
    </div>
  </header>
</template>

<style lang="scss" scoped>
header {
  padding: 2px;
  box-sizing: border-box;
  width: 100%;
  background: var(--header);
  display: flex;
  justify-content: flex-end;
  align-items: center;
  border-bottom: 2px solid var(--base);
  -webkit-app-region: drag;
  height: var(--header-height);
  z-index: 100;
  position: relative;
  .buttons {
    display: flex;
    gap: 2px;
  }
  .user-name {
    flex-grow: 1;
    text-align: left;
    padding: var(--padding-medium);
    padding-top: calc(var(--padding-medium) + 2px);
    box-sizing: border-box;
    font-size: 0.866rem;
    filter: opacity(0.5);
  }
  button {
    border: 2px solid white;
    padding: 0;
    margin: 0;
    filter: opacity(0.666);
    box-shadow: var(--shadow-medium);
    transition: filter 100ms ease;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--border-radius-medium);
    -webkit-app-region: no-drag;
    width: 1.066rem;
    height: 1.066rem;
  }
  button:hover {
    border: 2px solid white;
    filter: opacity(1);
  }
  .account-button svg {
    width: 80%;
    height: 80%;
  }
  .quit-button {
    background: red;
  }
  .minimize-button {
    svg {
      width: 70%;
      height: 70%;
    }
  }
}
.account-menu {
  display: flex;
  flex-direction: column;
  position: fixed;
  top: calc(var(--header-height) + 1px);
  right: 2px;
  background: var(--foreground);
  border: 2px solid var(--base);
  border-radius: var(--border-radius-medium);
  overflow: hidden;
  box-shadow: var(--shadow-medium);
  button {
    padding: 0;
    margin: 0;
    border: none;
    background: none;
    border-radius: 0;
    padding: 0 var(--padding-medium);
    display: flex;
    min-width: fit-content;
    width: 100%;
    justify-content: flex-end;
    box-shadow: none;
    font-size: 1rem;
  }
  .button,
  .button:hover {
    border-bottom: 1px solid var(--base-light);
  }
  .button:last-of-type,
  .button:last-of-type:hover {
    border-bottom: none;
  }
  .button:hover {
    border: none;
    background: var(--foreground-hover);
  }
}
</style>
