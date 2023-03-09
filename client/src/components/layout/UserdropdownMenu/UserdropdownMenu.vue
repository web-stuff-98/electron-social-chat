<script lang="ts" setup>
import { getOwnRoomIDs } from "../../../services/Rooms";
import { onMounted, watch, ref, onBeforeUnmount } from "vue";
import { socketStore } from "../../../store/SocketStore";
import { userdropdownStore } from "../../../store/UserDropdownStore";
import { IResMsg } from "../../../interfaces/GeneralInterfaces";
import ResMsg from "../ResMsg.vue";
import InviteToRoomCard from "../InviteToRoomCard.vue";

enum EUserdropdownMenuSection {
  "MENU" = "Menu",
  "INVITE_TO_ROOM" = "Invite to room",
  "DIRECT_MESSAGE" = "Direct message",
}

const mousePos = ref<{ left: number; top: number }>({ left: 0, top: 0 });
const menuPos = ref<{ left: number; top: number }>({ left: 0, top: 0 });
const mouseInside = ref(false);
const handleMouseMove = (e: MouseEvent) =>
  (mousePos.value = { left: e.clientX, top: e.clientY });
const section = ref<EUserdropdownMenuSection>(EUserdropdownMenuSection.MENU);
const handleMouseEnter = () => (mouseInside.value = true);
const handleMouseLeave = () => (mouseInside.value = false);
const getOwnRoomIDsResMsg = ref<IResMsg>({ msg: "", err: false, pen: false });

watch(userdropdownStore, () => {
  menuPos.value = mousePos.value;
  section.value = EUserdropdownMenuSection.MENU;
  getOwnRoomIDsResMsg.value = { msg: "", err: false, pen: false };
});

onMounted(() => {
  window.addEventListener("mousemove", handleMouseMove);
});

onBeforeUnmount(() => {
  window.removeEventListener("mousemove", handleMouseMove);
});

const directMessageClicked = () =>
  (section.value = EUserdropdownMenuSection.DIRECT_MESSAGE);

const ownRoomIDs = ref<string[]>([]);
async function inviteToRoomClicked() {
  section.value = EUserdropdownMenuSection.INVITE_TO_ROOM;
  try {
    ownRoomIDs.value = [];
    getOwnRoomIDsResMsg.value = { msg: "", err: false, pen: true };
    const ids: string[] = await getOwnRoomIDs();
    ownRoomIDs.value = ids;
    getOwnRoomIDsResMsg.value = { msg: ids.length > 0 ? "" : "You have no rooms", err: false, pen: false };
  } catch (e) {
    getOwnRoomIDsResMsg.value = { msg: `${e}`, err: true, pen: false };
  }
}

function inviteToRoom(roomId: string) {
  socketStore.send(
    JSON.stringify({
      event_type: "ROOM_INVITATION",
      recipient: userdropdownStore.subject,
      room_id: roomId,
    })
  );
  userdropdownStore.open = false;
}

function friendRequestClicked() {
  socketStore.send(
    JSON.stringify({
      event_type: "FRIEND_REQUEST",
      recipient: userdropdownStore.subject,
    })
  );
  userdropdownStore.open = false;
}

function blockClicked() {
  socketStore.send(
    JSON.stringify({
      event_type: "BLOCK_USER",
      uid: userdropdownStore.subject,
    })
  );
  userdropdownStore.open = false;
}

function banClicked() {
  socketStore.send(
    JSON.stringify({
      event_type: "BAN_USER",
      uid: userdropdownStore.subject,
      room_id: userdropdownStore.roomId,
    })
  );
  userdropdownStore.open = false;
}

const msgInputRef = ref<HTMLCanvasElement | null>();
const msgInput = ref("");
function handleMsgInput(e: Event) {
  const target = e.target as HTMLInputElement;
  if (!target || !target.value || target.value.length > 300) return;
  msgInput.value = target.value;
}
function submitDirectMessage() {
  if (msgInput.value == "" || msgInput.value.length > 300) return;
  socketStore.send(
    JSON.stringify({
      event_type: "DIRECT_MESSAGE",
      content: msgInput.value,
      recipient: userdropdownStore.subject,
    })
  );
  // @ts-ignore
  msgInputRef.value = "";
  msgInput.value = "";
  userdropdownStore.open = false;
}
</script>

<template>
  <div
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
    v-show="userdropdownStore.open"
    class="container"
    :style="{ left: `${menuPos.left}px`, top: `${menuPos.top}px` }"
  >
    <!-- Menu section -->
    <div v-if="section === EUserdropdownMenuSection.MENU" class="menu">
      <button @click="inviteToRoomClicked">Invite to room</button>
      <button @click="directMessageClicked">Direct message</button>
      <button @click="friendRequestClicked">Friend request</button>
      <button @click="blockClicked">Block</button>
      <button v-if="userdropdownStore.roomId" @click="banClicked">Ban</button>
    </div>
    <!-- Direct message section -->
    <form
      @submit.prevent="submitDirectMessage"
      v-if="section === EUserdropdownMenuSection.DIRECT_MESSAGE"
      class="direct-message"
    >
      <input maxlength="300" @input="handleMsgInput" ref="msgInputRef" />
      <button type="submit">
        <v-icon name="md-send" />
      </button>
    </form>
    <!-- Invite to room section -->
    <div
      v-if="section === EUserdropdownMenuSection.INVITE_TO_ROOM"
      class="invite-to-room"
    >
      <ResMsg :resMsg="getOwnRoomIDsResMsg" />
      <div
        @click="() => inviteToRoom(id)"
        class="room-container"
        v-for="id in ownRoomIDs"
      >
        <InviteToRoomCard :id="id" />
      </div>
    </div>
    <!-- Close button -->
    <button @click="userdropdownStore.open = false" class="close-button">
      <v-icon name="io-close" />
    </button>
  </div>
</template>

<style lang="scss" scoped>
.container {
  position: fixed;
  padding: 2px;
  gap: 2px;
  background: var(--foreground);
  border: 1px solid var(--base-light);
  box-shadow: var(--shadow);
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  border-radius: var(--border-radius-medium);
  z-index: 100;
  .menu {
    padding: 0;
    width: fit-content;
    display: flex;
    gap: 2px;
    flex-direction: column;
    align-items: center;
    button {
      display: flex;
    }
  }
  .direct-message {
    display: flex;
    gap: 2px;
    display: flex;
    align-items: center;
    button {
      display: flex;
      border: none;
      box-shadow: none;
      background: none;
    }
  }
  .invite-to-room {
    display: flex;
    flex-direction: column;
    gap: 2px;
    .room-container {
      padding: 0;
    }
  }
  button {
    padding: var(--padding-medium);
    text-align: left;
    box-shadow: none;
    flex-grow: 1;
    width: 100%;
  }
  .close-button {
    border: 1px solid var(--text-color);
    padding: 0;
    filter: opacity(0.666);
    width: fit-content;
    box-shadow: var(--shadow-medium);
    transition: filter 100ms ease;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--border-radius-medium);
    background: red;
    svg {
      width: 1rem;
      height: 1rem;
    }
  }
  .close-button:hover {
    outline: 1px solid var(--text-color);
    filter: opacity(1);
  }
}
</style>
