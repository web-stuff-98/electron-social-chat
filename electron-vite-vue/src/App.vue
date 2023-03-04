<script setup lang="ts">
import WelcomeModal from "./components/welcomeModal/WelcomeModal.vue";
import Bar from "./components/layout/Bar.vue";
import AsideMenu from "./components/layout/asideMenu/AsideMenu.vue";

import { authStore, IUser } from "./store/AuthStore";
import { socketStore } from "./store/SocketStore";
import { userStore } from "./store/UserStore";

import { useRouter } from "vue-router";
import { watch, onBeforeUnmount, onMounted, ref } from "vue";
import {
  parseSocketEventData,
  instanceOfChangeData,
} from "./utils/determineSocketEvent";
import { roomStore } from "./store/RoomStore";
import { baseURL } from "./services/makeRequest";

const router = useRouter();
const showAside = ref(false);

/* ------- Update users when socket event received ------- */
const watchForUserChanges = (e: MessageEvent) => {
  const data = parseSocketEventData(e);
  if (!data) return;
  if (instanceOfChangeData(data)) {
    if (data.ENTITY === "USER") {
      if (data.METHOD === "UPDATE" || "UPDATE_IMAGE") {
        if (data.DATA.ID === authStore.user?.ID) {
          authStore.user = {
            ...authStore.user,
            ...(data.DATA as Partial<IUser>),
          };
        } else {
          const i = userStore.users.findIndex((u) => u.ID === data.DATA.ID);
          userStore.users[i] = { ...userStore.users[i], ...data.DATA };
        }
      }
    }
  }
};

/* ------- Update rooms when socket event received ------- */
const watchForRoomChanges = (e: MessageEvent) => {
  const data = parseSocketEventData(e);
  if (!data) return;
  if (instanceOfChangeData(data)) {
    if (data.ENTITY === "ROOM") {
      if (data.METHOD === "UPDATE") {
        const i = roomStore.rooms.findIndex((r) => r.ID === data.DATA.ID);
        roomStore.rooms[i] = { ...roomStore.rooms[i], ...data.DATA };
      }
      if (data.METHOD === "DELETE") {
        const i = roomStore.rooms.findIndex((r) => r.ID === data.DATA.ID);
        roomStore.rooms.splice(i, 1);
        if (roomStore.currentRoom === data.DATA.ID) {
          router.push("/");
        }
      }
      if (data.METHOD === "UPDATE_IMAGE") {
        console.log("UPDATE IMAGE");
        const i = roomStore.rooms.findIndex((r) => r.ID === data.DATA.ID);
        let imgUrl = roomStore.rooms[i].img_url || "";
        if (imgUrl) {
          const split = imgUrl.split("?v=");
          imgUrl = split[0] + `?v=${Number(split[1]) + 1}`;
        } else {
          imgUrl = `${baseURL}/api/room/image/${data.DATA.ID}?v=1`;
        }
        roomStore.rooms[i].img_url = imgUrl;
      }
    }
  }
};

watch(socketStore, (_, newVal) => {
  if (newVal.socket) {
    socketStore.socket?.addEventListener("message", watchForUserChanges);
    socketStore.socket?.addEventListener("message", watchForRoomChanges);
  }
});

const clearUserCacheInterval = ref<NodeJS.Timer>();
const clearRoomCacheInterval = ref<NodeJS.Timer>();
const refreshTokenInterval = ref<NodeJS.Timer>();

onBeforeUnmount(() => {
  /* ------- Cleanup intervals ------- */
  clearInterval(clearUserCacheInterval.value);
  clearInterval(clearRoomCacheInterval.value);
  clearInterval(refreshTokenInterval.value);

  /* ------- Cleanup socket event listeners ------- */
  socketStore.socket?.removeEventListener("message", watchForUserChanges);
  socketStore.socket?.removeEventListener("message", watchForRoomChanges);
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
  /* ------- Clear rooms in cache interval ------- */
  clearRoomCacheInterval.value = setInterval(() => {
    roomStore.disappearedRooms.forEach((r) => {
      if (Date.now() > r.lastSeen + 30000) {
        const found = roomStore.visibleRooms.find((id) => id === r.ID);
        if (!found) {
          roomStore.rooms = roomStore.rooms.filter((r) => r.ID !== found);
          roomStore.disappearedRooms = roomStore.disappearedRooms.filter(
            (r) => r.ID !== found
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
