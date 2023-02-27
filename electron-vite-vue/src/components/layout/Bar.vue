<script lang="ts">
import { ipcRenderer } from "electron";
import { ref } from "vue";
import { authStore } from "../../store/AuthStore";
import { modalStore, EModalType } from "../../store/ModalStore";

export default {
  setup() {
    const showAccountMenu = ref(false);
    const mouseInAccountMenu = ref(false);
    return {
      authStore,
      modalStore,

      showAccountMenu,
      mouseInAccountMenu,
      
      methods: {
        quit() {
          ipcRenderer.send("window", ["QUIT"]);
        },
        minimize() {
          ipcRenderer.send("window", ["MINIMIZE"]);
        },
        click() {
          if (!mouseInAccountMenu.value) showAccountMenu.value = false;
        },
        deleteAccount() {
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
        },
      },
    };
  },
  mounted() {
    window.document.addEventListener("click", this.methods.click);
  },
  beforeDestroy() {
    window.document.removeEventListener("click", this.methods.click);
  },
};
</script>

<template>
  <header>
    <div v-if="authStore.user" class="user-name">
      {{ authStore.user?.username }}
    </div>
    <div class="header-buttons">
      <div
        @mouseenter="mouseInAccountMenu = true"
        @mouseleave="mouseInAccountMenu = false"
        v-show="showAccountMenu && !modalStore.showModal"
        class="account-menu"
      >
        <button @click="methods.deleteAccount">Delete account</button>
        <button @click="authStore.logout">Log out</button>
      </div>
      <button
        @mouseenter="mouseInAccountMenu = true"
        @mouseleave="mouseInAccountMenu = false"
        @click="showAccountMenu = !showAccountMenu"
        class="account-button"
      >
        <v-icon name="md-manageaccounts-sharp" />
      </button>
      <button @click="methods.minimize" class="minimize-button">
        <v-icon name="fa-minus" />
      </button>
      <button @click="methods.quit" class="quit-button">
        <v-icon name="io-close" />
      </button>
    </div>
  </header>
</template>

<style scoped>
header {
  padding: 3px;
  box-sizing: border-box;
  width: 100%;
  background: var(--header);
  display: flex;
  justify-content: flex-end;
  align-items: center;
  border-bottom: 1px solid var(--base);
  -webkit-app-region: drag;
  height: var(--header-height);
  z-index: 100;
  position: relative;
}
.user-name {
  flex-grow: 1;
  text-align: left;
  padding: var(--padding-medium);
  padding-top: calc(var(--padding-medium) + 2px);
  box-sizing: border-box;
  font-size: 1.125rem;
}
.header-buttons {
  display: flex;
  gap: 2px;
}
header button {
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
}
header button:hover {
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
header button {
  width: 1.5rem;
  height: 1.5rem;
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
}
.account-menu button {
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
.account-menu button,
.account-menu button:hover {
  border-bottom: 1px solid var(--base-light);
}
.account-menu button:last-of-type,
.account-menu button:last-of-type:hover {
  border-bottom: none;
}
.account-menu button:hover {
  border: none;
  border-bottom: 1px solid var(--base-light);
  background: var(--foreground-hover);
}
</style>
