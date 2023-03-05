<script setup lang="ts">
import { ref } from "vue";
import { authStore } from "../../store/AuthStore";

import { onBeforeMount, onBeforeUnmount } from "vue";
import MessageModal from "../messageModal/MessageModal.vue";
import { IResMsg } from "../../interfaces/GeneralInterfaces";

const showAccountMenu = ref(false);
const mouseInAccountMenu = ref(false);

const modalConfirmation = ref(() => {});
const modalCancellation = ref<Function | undefined>(() => {});
const showModal = ref(false);
const modalMsg = ref<IResMsg>({ msg: "", err: false, pen: false });

function quit() {
  // quit
}
function minimize() {
  // minimize
}
function click() {
  if (!mouseInAccountMenu.value) showAccountMenu.value = false;
}
function deleteAccount() {
  modalMsg.value = {
    msg: "Are you sure you want to delete your account?",
    err: false,
    pen: false,
  };
  showModal.value = true;
  modalConfirmation.value = () => {
    authStore.deleteAccount();
    showModal.value = false;
  };
  showAccountMenu.value = false;
  modalCancellation.value = () => (showModal.value = false);
}

onBeforeMount(() => {
  window.document.addEventListener("click", click);
});
onBeforeUnmount(() => {
  window.document.removeEventListener("click", click);
});
</script>

<template>
  <MessageModal
    :msg="modalMsg"
    :show="showModal"
    :confirmationCallback="modalConfirmation"
    :cancellationCallback="modalCancellation"
  />
  <header>
    <div v-if="authStore.user" class="user-name">
      {{ authStore.user?.username }}
    </div>
    <div class="buttons">
      <div
        @mouseenter="mouseInAccountMenu = true"
        @mouseleave="mouseInAccountMenu = false"
        v-show="showAccountMenu && !showModal"
        class="account-menu"
      >
        <button @click="deleteAccount">Delete account</button>
        <button
          @click="
            {
              authStore.logout();
              showAccountMenu = false;
            }
          "
        >
          Log out
        </button>
      </div>
      <button
        v-if="authStore.user"
        @mouseenter="mouseInAccountMenu = true"
        @mouseleave="mouseInAccountMenu = false"
        @click="showAccountMenu = !showAccountMenu"
        class="account-button"
      >
        <v-icon name="bi-shield-fill" />
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
  padding: 3px;
  padding-right: 2px;
  box-sizing: border-box;
  width: 100%;
  background: var(--header);
  display: flex;
  justify-content: flex-end;
  align-items: center;
  border-bottom: 2px solid var(--base-light);
  -webkit-app-region: drag;
  height: var(--header-height);
  z-index: 100;
  position: fixed;
  top: 0;
  .buttons {
    display: flex;
    gap: 1px;
  }
  .user-name {
    flex-grow: 1;
    text-align: left;
    padding: var(--padding-medium);
    padding-top: calc(var(--padding-medium) + 4px);
    box-sizing: border-box;
    font-size: 0.866rem;
    filter: opacity(0.5);
  }
  button {
    border: 2px solid lightgray;
    padding: 0;
    margin: 0;
    filter: opacity(0.666);
    box-shadow: var(--shadow-medium);
    transition: filter 100ms ease;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 33.33%;
    -webkit-app-region: no-drag;
    width: 1.3rem;
    height: 1.3rem;
    background: rgb(38, 38, 38);
  }
  button:hover {
    border: 2px solid white;
    filter: opacity(1);
  }
  .account-button svg {
    width: 90%;
    height: 90%;
  }
  .quit-button {
    background: red;
    svg {
      width: 135%;
      height: 135%;
    }
  }
  .minimize-button {
    svg {
      width: 90%;
      height: 90%;
    }
  }
}
.account-menu {
  display: flex;
  flex-direction: column;
  position: fixed;
  top: calc(var(--header-height) + 1px);
  padding: 2px;
  right: 2px;
  background: var(--foreground);
  border: 1px solid var(--base-light);
  border-radius: var(--border-radius-medium);
  overflow: hidden;
  box-shadow: var(--shadow-medium);
  button {
    margin: 0;
    border: none;
    background: none;
    border-radius: 0;
    padding: var(--padding-medium);
    display: flex;
    min-width: fit-content;
    width: 100%;
    justify-content: flex-end;
    box-shadow: none;
    font-size: 1rem;
  }
  button:hover {
    border: none;
    background: var(--foreground-hover);
    border-radius: 3px;
  }
}
</style>
