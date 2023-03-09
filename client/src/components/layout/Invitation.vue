<script lang="ts" setup>
import { onMounted, onBeforeUnmount, toRefs } from "vue";
import { IInvitation } from "../../interfaces/GeneralInterfaces";
import { authStore } from "../../store/AuthStore";
import { roomStore } from "../../store/RoomStore";
import { socketStore } from "../../store/SocketStore";
import { userStore } from "../../store/UserStore";

const props = defineProps<{ inv: IInvitation }>();
const { inv } = toRefs(props);

onMounted(() => {
  roomStore.roomEnteredView(inv.value.room_id);
  userStore.userEnteredView(inv.value.author);
});
onBeforeUnmount(() => {
  roomStore.roomLeftView(inv.value.room_id);
  userStore.userEnteredView(inv.value.author);
});

function accept() {
  socketStore.send(
    JSON.stringify({
      event_type: "ROOM_INVITATION_RESPONSE",
      ID: inv.value.ID,
      accept: true,
    })
  );
}
function decline() {
  socketStore.send(
    JSON.stringify({
      event_type: "ROOM_INVITATION_RESPONSE",
      ID: inv.value.ID,
      accept: false,
    })
  );
}
</script>

<template>
  <div class="container">
    <div
      v-if="!inv.accepted && !inv.declined && inv.author !== authStore.user?.ID"
      class="message"
    >
      {{ userStore.getUser(inv.author)?.username + " " }}
      invited you to join
      {{ " " + roomStore.getRoom(inv.room_id)?.name }}
    </div>
    <div
      v-if="!inv.accepted && !inv.declined && inv.author === authStore.user?.ID"
      class="message"
    >
      You sent an invitation to
      {{ " " + userStore.getUser(inv.recipient)?.username + " " }}
      to join
      {{ " " + roomStore.getRoom(inv.room_id)?.name }}
    </div>
    <div
      v-if="inv.accepted && inv.author === authStore.user?.ID"
      class="message"
    >
      {{ userStore.getUser(inv.recipient)?.username + " " }}
      accepted your invitation
    </div>
    <div
      v-if="inv.accepted && inv.author !== authStore.user?.ID"
      class="message"
    >
      You accepted
      {{ " " + userStore.getUser(inv.author)?.username + "'s " }}
      invitation
    </div>
    <div
      v-if="inv.declined && inv.author === authStore.user?.ID"
      class="message"
    >
      {{ userStore.getUser(inv.recipient)?.username + " " }}
      declined your invitation
    </div>
    <div
      v-if="inv.declined && inv.author !== authStore.user?.ID"
      class="message"
    >
      You declined
      {{ " " + userStore.getUser(inv.author)?.username + "'s " }}
      invitation
    </div>
    <div
      v-if="inv.author !== authStore.user?.ID && !inv.accepted && !inv.declined"
      class="buttons"
    >
      <button type="button" @click="accept">Accept</button>
      <button type="button" @click="decline">Decline</button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  gap: var(--padding-medium);
  padding: var(--padding-medium);
  box-sizing: border-box;
  .message {
    text-align: center;
  }
  .buttons {
    width: 100%;
    display: flex;
    gap: var(--padding-medium);
    button {
      width: 50%;
    }
  }
}
</style>
