<script setup lang="ts">
import Modal from "./components/modal/Modal.vue";
import Bar from "./components/layout/Bar.vue";
import AsideMenu from "./components/layout/asideMenu/AsideMenu.vue";

import { authStore, IUser } from "./store/AuthStore";
import { modalStore } from "./store/ModalStore";
import { socketStore } from "./store/SocketStore";
import { userStore } from "./store/UserStore";

import { watch, onBeforeUnmount, onMounted, ref } from "vue";

import {
  parseSocketEventData,
  instanceOfChangeData,
} from "./utils/determineSocketEvent";

/* ------- Connect socket when user assigned (resend cookie) ------- */
watch(authStore, (_, newVal) => {
  if (newVal.user) {
    socketStore.socket = new WebSocket(
      process.env.NODE_ENV === "development" ||
      window.location.origin === "http://localhost:8080"
        ? "ws://localhost:8080/api/ws"
        : "wss://electron-social-chat-backend.herokuapp.com/api/ws"
    );
    socketStore.openSubscription(`user=${newVal.user.ID}`);
  } else {
    socketStore.socket = undefined;
    socketStore.subscriptions = [];
  }
});

/* ------- Update current user in authStore when socket event received ------- */
const watchForCurrentUserChanges = (e: MessageEvent) => {
  console.log("Message received:", e.data);
  const data = parseSocketEventData(e);
  if (!data) return;
  if (instanceOfChangeData(data)) {
    if (data.DATA.ID === authStore.user?.ID) {
      authStore.user = {
        ...authStore.user,
        ...(data.DATA as Partial<IUser>),
      };
    }
  }
};
watch(socketStore, (_, newVal) => {
  if (newVal.socket) {
    socketStore.socket?.addEventListener("message", watchForCurrentUserChanges);
  }
});

const clearUserCacheInterval = ref<NodeJS.Timer>();
const refreshTokenInterval = ref<NodeJS.Timer>();

onBeforeUnmount(() => {
  /* ------- Cleanup intervals ------- */
  clearInterval(clearUserCacheInterval.value);
  clearInterval(refreshTokenInterval.value);

  /* ------- Cleanup socket event listeners ------- */
  socketStore.socket?.removeEventListener(
    "message",
    watchForCurrentUserChanges
  );
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
  <Bar />
  <AsideMenu v-if="authStore.user" />
  <Modal v-if="modalStore.showModal || !authStore.user" />
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
