<script setup lang="ts">
import WelcomeModal from "./components/welcomeModal/WelcomeModal.vue";
import AsideMenu from "./components/layout/asideMenu/AsideMenu.vue";

import { authStore, IUser } from "./store/AuthStore";
import { socketStore } from "./store/SocketStore";
import { userStore } from "./store/UserStore";

import { useRouter } from "vue-router";
import { watch, onBeforeUnmount, onMounted, ref } from "vue";
import {
  parseSocketEventData,
  instanceOfChangeData,
  instanceOfResponseMessageData,
  instanceOfDirectMessageData,
  instanceOfDirectMessageUpdateData,
  instanceOfDirectMessageDeleteData,
  instanceOfRoomInvitationData,
  instanceOfRoomInvitationDeleteData,
  instanceOfRoomInvitationResponseData,
  instanceOfFriendRequestData,
  instanceOfFriendRequestDeleteData,
  instanceOfFriendRequestResponseData,
} from "./utils/determineSocketEvent";
import { roomStore } from "./store/RoomStore";
import { baseURL } from "./services/makeRequest";
import MessageModal from "./components/messageModal/MessageModal.vue";
import { IResMsg } from "./interfaces/GeneralInterfaces";
import UserdropdownMenu from "./components/layout/UserdropdownMenu/UserdropdownMenu.vue";
import { messagingStore } from "./store/MessagingStore";

const router = useRouter();
const showAside = ref(false);

const modalConfirmation = ref(() => {});
const modalCancellation = ref<Function | undefined>(() => {});
const showModal = ref(false);
const modalMsg = ref<IResMsg>({ msg: "", err: false, pen: false });

/* ------- Update users when socket event received ------- */
const watchForUserChanges = (e: MessageEvent) => {
  const data = parseSocketEventData(e);
  if (!data) return;
  if (instanceOfChangeData(data)) {
    if (data.ENTITY === "USER") {
      if (data.METHOD === "UPDATE") {
        if (data.DATA.ID === authStore.user?.ID) {
          authStore.user = {
            ...authStore.user,
            ...(data.DATA as Partial<IUser>),
          };
        } else {
          const i = userStore.users.findIndex((u) => u.ID === data.DATA.ID);
          userStore.users[i] = { ...userStore.users[i], ...data.DATA };
        }
      } else if (data.METHOD === "UPDATE_IMAGE") {
        userStore.cacheUserData(data.DATA.ID, true);
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

/* ------- Watch for direct messages, invites & friend requests on the socket connection ------- */
const watchMessaging = (e: MessageEvent) => {
  const data = parseSocketEventData(e);
  if (!data) return;
  if (instanceOfDirectMessageData(data)) {
    const convI = messagingStore.conversations.findIndex(
      (c) =>
        c.uid ===
        (data.author === authStore.user?.ID ? data.recipient : data.author)
    );
    console.log(convI);
    if (convI !== -1)
      messagingStore.conversations[convI].messages.push({
        ...data,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      });
    else
      messagingStore.conversations.push({
        uid: data.author === authStore.user?.ID ? data.recipient : data.author,
        messages: [
          {
            ...data,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
          },
        ],
        friend_requests: [],
        invitations: [],
      });
  }
  if (instanceOfDirectMessageUpdateData(data)) {
    const convI = messagingStore.conversations.findIndex(
      (c) =>
        c.uid ===
        (data.author === authStore.user?.ID ? data.recipient : data.author)
    );
    const msgI = messagingStore.conversations[convI].messages.findIndex(
      (m) => m.ID === data.ID
    );
    messagingStore.conversations[convI].messages[msgI].content = data.content;
    messagingStore.conversations[convI].messages[msgI].updated_at =
      new Date().toISOString();
  }
  if (instanceOfDirectMessageDeleteData(data)) {
    const convI = messagingStore.conversations.findIndex(
      (c) =>
        c.uid ===
        (data.author === authStore.user?.ID ? data.recipient : data.author)
    );
    const msgI = messagingStore.conversations[convI].messages.findIndex(
      (m) => m.ID === data.ID
    );
    messagingStore.conversations[convI].messages.splice(msgI, 1);
  }
  if (instanceOfRoomInvitationData(data)) {
    const convI = messagingStore.conversations.findIndex(
      (c) =>
        c.uid ===
        (data.author === authStore.user?.ID ? data.recipient : data.author)
    );
    if (convI !== -1)
      messagingStore.conversations[convI].invitations = [
        ...messagingStore.conversations[convI].invitations,
        {
          ...data,
          accepted: false,
          declined: false,
          created_at: new Date().toISOString(),
        },
      ];
    else
      messagingStore.conversations.push({
        uid: data.author === authStore.user?.ID ? data.recipient : data.author,
        messages: [],
        invitations: [
          {
            ...data,
            accepted: false,
            declined: false,
            created_at: new Date().toISOString(),
          },
        ],
        friend_requests: [],
      });
  }
  if (instanceOfRoomInvitationDeleteData(data)) {
    const convI = messagingStore.conversations.findIndex(
      (c) =>
        c.uid ===
        (data.author === authStore.user?.ID ? data.recipient : data.author)
    );
    const invI = messagingStore.conversations[convI].invitations.findIndex(
      (r) => r.ID === data.ID
    );
    messagingStore.conversations[convI].invitations.splice(invI, 1);
  }
  if (instanceOfRoomInvitationResponseData(data)) {
    const convI = messagingStore.conversations.findIndex(
      (c) =>
        c.uid ===
        (data.author === authStore.user?.ID ? data.recipient : data.author)
    );
    const invI = messagingStore.conversations[convI].friend_requests.findIndex(
      (r) => r.ID === data.ID
    );
    messagingStore.conversations[convI].invitations[invI].accepted =
      data.accept;
    messagingStore.conversations[convI].invitations[invI].declined =
      !data.accept;
  }
};

/* ------- Watch for response messages on the socket connection ------- */
const watchForResponseMessages = (e: MessageEvent) => {
  const data = parseSocketEventData(e);
  if (!data) return;
  if (instanceOfResponseMessageData(data)) {
    modalMsg.value = {
      msg: data.DATA.msg,
      err: data.DATA.err,
      pen: false,
    };
    modalConfirmation.value = () => (showModal.value = false);
    modalCancellation.value = undefined;
    showModal.value = true;
  }
};

watch(socketStore, (_, newVal) => {
  if (newVal.socket) {
    socketStore.socket?.addEventListener("message", watchForUserChanges);
    socketStore.socket?.addEventListener("message", watchForRoomChanges);
    socketStore.socket?.addEventListener("message", watchForResponseMessages);
    socketStore.socket?.addEventListener("message", watchMessaging);
  } else {
    socketStore.connectSocket(authStore.user?.ID!);
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
  socketStore.socket?.removeEventListener("message", watchForResponseMessages);
  socketStore.socket?.removeEventListener("message", watchMessaging);
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
    <MessageModal
      :msg="modalMsg"
      :show="showModal"
      :confirmationCallback="modalConfirmation"
      :cancellationCallback="modalCancellation"
    />
    <UserdropdownMenu v-if="authStore.user" />
    <main :style="showAside ? { 'padding-left': 'var(--aside-width)' } : {}">
      <router-view />
    </main>
    <AsideMenu
      :toggleShowAside="() => (showAside = !showAside)"
      :showAside="showAside"
      v-if="authStore.user"
    />
    <WelcomeModal v-else />
  </div>
</template>

<style scoped lang="scss">
.root {
  display: flex;
  height: 100%;
  box-sizing: border-box;
  background: var(--background-radial);

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
