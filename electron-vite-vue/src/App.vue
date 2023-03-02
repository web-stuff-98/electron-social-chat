<script setup lang="ts">
import WelcomeModal from "./components/welcomeModal/WelcomeModal.vue";
import Bar from "./components/layout/Bar.vue";
import AsideMenu from "./components/layout/asideMenu/AsideMenu.vue";

import { authStore, IUser } from "./store/AuthStore";
import { socketStore } from "./store/SocketStore";
import { userStore } from "./store/UserStore";

import { watch, onBeforeUnmount, onMounted, ref } from "vue";

import {
  parseSocketEventData,
  instanceOfChangeData,
} from "./utils/determineSocketEvent";

const showAside = ref(false);

/* ------- Update current user in authStore when socket event received ------- */
const watchForCurrentUserChanges = (e: MessageEvent) => {
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
  <div class="root">
    <Bar />
    <AsideMenu
      :toggleShowAside="() => (showAside = !showAside)"
      :showAside="showAside"
      v-if="authStore.user"
    />
    <WelcomeModal v-else />
    <main :style="showAside ? { 'padding-left': 'var(--aside-width)' } : {}">
      <router-view />
    </main>
  </div>
</template>

<style scoped lang="scss">
.root {
  display: flex;
  height: 100%;
  padding-top: var(--header-height);
  box-sizing: border-box;
  main {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
    flex-grow: 1;
    flex-direction: column;
  }
}
</style>
