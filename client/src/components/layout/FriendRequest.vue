<script lang="ts" setup>
import { onMounted, onBeforeUnmount, toRefs } from "vue";
import { IFriendRequest } from "../../interfaces/GeneralInterfaces";
import { authStore } from "../../store/AuthStore";
import { socketStore } from "../../store/SocketStore";
import { userStore } from "../../store/UserStore";

const props = defineProps<{ frq: IFriendRequest }>();
const { frq } = toRefs(props);

onMounted(() => {
  userStore.userEnteredView(frq.value.author);
});
onBeforeUnmount(() => {
  userStore.userEnteredView(frq.value.author);
});

function accept() {
  socketStore.send(
    JSON.stringify({
      event_type: "FRIEND_REQUEST_RESPONSE",
      ID: frq.value.ID,
      accept: true,
    })
  );
}
function decline() {
  socketStore.send(
    JSON.stringify({
      event_type: "FRIEND_REQUEST_RESPONSE",
      ID: frq.value.ID,
      accept: false,
    })
  );
}
</script>

<template>
  <div class="request">
    <div
      v-if="!frq.accepted && !frq.declined && frq.author !== authStore.user?.ID"
      class="message"
    >
      {{ userStore.getUser(frq.author)?.username + " " }} sent you a friend
      request
    </div>
    <div
      v-if="!frq.accepted && !frq.declined && frq.author === authStore.user?.ID"
      class="message"
    >
      You sent a friend request to
      {{ " " + userStore.getUser(frq.recipient)?.username }}
    </div>
    <div
      v-if="frq.accepted && frq.author === authStore.user?.ID"
      class="message"
    >
      {{ userStore.getUser(frq.recipient)?.username + " " }}
      accepted your friend request
    </div>
    <div
      v-if="frq.accepted && frq.author !== authStore.user?.ID"
      class="message"
    >
      You accepted
      {{ " " + userStore.getUser(frq.author)?.username + "'s " }}
      friend request
    </div>
    <div
      v-if="frq.declined && frq.author === authStore.user?.ID"
      class="message"
    >
      {{ userStore.getUser(frq.recipient)?.username + " " }}
      declined your friend request
    </div>
    <div
      v-if="frq.declined && frq.author !== authStore.user?.ID"
      class="message"
    >
      You declined
      {{ " " + userStore.getUser(frq.author)?.username + "'s " }}
      friend request
    </div>
    <div
      v-if="frq.author !== authStore.user?.ID && !frq.accepted && !frq.declined"
      class="buttons"
    >
      <button type="button" @click="accept">Accept</button>
      <button type="button" @click="decline">Decline</button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.request {
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
    font-size: 0.7rem;
    box-sizing: border-box;
    font-weight: 600;
  }
  .buttons {
    width: 100%;
    display: flex;
    gap: var(--padding-medium);
    box-sizing: border-box;
    button {
      padding: 3px;
      width: 50%;
      font-size: 0.7rem;
      border: 2px solid var(--base-light);
    }
  }
}
</style>
