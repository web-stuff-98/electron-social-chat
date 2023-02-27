<script setup lang="ts">
import Modal from "./components/modal/Modal.vue";
import Bar from "./components/layout/Bar.vue";
import AsideMenu from "./components/layout/asideMenu/AsideMenu.vue";

import { authStore } from "./store/AuthStore";
import { modalStore } from "./store/ModalStore";
import { socketStore } from "./store/SocketStore";
import { userStore } from "./store/UserStore";

import { watch, onBeforeUnmount, onMounted, ref } from "vue";

/* ------- Connect socket when user assigned (resend cookie) ------- */
watch(authStore, (_, newVal) => {
  if (newVal.user) {
    socketStore.socket = new WebSocket(
      process.env.NODE_ENV === "development" ||
      window.location.origin === "http://localhost:8080"
        ? "ws://localhost:8080/api/ws"
        : "wss://electron-social-chat-backend.herokuapp.com/api/ws"
    );
  } else {
    socketStore.socket = undefined;
  }
});

const clearUserCacheInterval = ref<NodeJS.Timer>();
const refreshTokenInterval = ref<NodeJS.Timer>();

onBeforeUnmount(() => {
  /* ------- Cleanup intervals ------- */
  clearInterval(clearUserCacheInterval.value);
  clearInterval(refreshTokenInterval.value);
});

onMounted(() => {
  /* ------- Clear users in cache interval ------- */
  clearUserCacheInterval.value = setInterval(() => {
    userStore.disappearedUsers.forEach((u) => {
      if (Date.now() > u.lastSeen + 30000) {
        const found = userStore.visibleUsers.find((uid) => uid === u.uid);
        if (!found) {
          userStore.users = userStore.users.filter((u) => u.ID !== found);
          userStore.disappearedUsers = userStore.disappearedUsers.filter(
            (u) => u.uid !== found
          );
        }
      }
    });
  }, 5000);

  /* ------- Refresh token interval ------- */
  refreshTokenInterval.value = setInterval(authStore.refreshToken, 100000);
});
</script>

<template>
  <Modal v-if="modalStore.showModal || !authStore.user" />
  <Bar />
  <AsideMenu v-if="authStore.user" />
  <div class="content">
    <main></main>
  </div>
</template>

<style>
.content {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
