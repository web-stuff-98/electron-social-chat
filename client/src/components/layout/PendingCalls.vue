<script lang="ts" setup>
import { authStore } from "../../store/AuthStore";
import { pendingCallsStore } from "../../store/CallsStore";
import { socketStore } from "../../store/SocketStore";
import { userStore } from "../../store/UserStore";

function cancelHangupClicked(index: number) {
  socketStore.send(
    JSON.stringify({
      event_type: "CALL_USER_RESPONSE",
      caller: pendingCallsStore[index].caller,
      called: pendingCallsStore[index].called,
      accept: false,
    })
  );
}

function acceptClicked(index: number) {
  socketStore.send(
    JSON.stringify({
      event_type: "CALL_USER_RESPONSE",
      caller: pendingCallsStore[index].caller,
      called: pendingCallsStore[index].called,
      accept: true,
    })
  );
}
</script>

<template>
  <div class="container">
    <div
      :style="{
        backgroundImage: `url(${userStore.getUser(
          pendingCall.caller === authStore.user?.ID
            ? pendingCall.called
            : pendingCall.caller
        )?.base64pfp})`,
      }"
      class="pending-call"
      v-for="(pendingCall, index) in pendingCallsStore"
    >
      {{
        !userStore.getUser(
          pendingCall.caller === authStore.user?.ID
            ? pendingCall.called
            : pendingCall.caller
        )?.base64pfp
          ? userStore.getUser(
              pendingCall.caller === authStore.user?.ID
                ? pendingCall.called
                : pendingCall.caller
            )?.username
          : ""
      }}
      <!-- Accept call button -->
      <button
        v-if="pendingCall.caller !== authStore.user?.ID"
        @click="() => acceptClicked(index)"
        class="accept-button"
      >
        <v-icon name="hi-phone-incoming" />
      </button>
      <!-- Cancel/Hangup call button -->
      <button
        @click="() => cancelHangupClicked(index)"
        class="cancel-hangup-button"
      >
        <v-icon name="hi-phone-missed-call" />
      </button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.container {
  position: fixed;
  bottom: 0;
  right: 0;
  display: flex;
  gap: 0.5rem;
  padding: 0.5rem;
  .pending-call {
    width: 4rem;
    height: 4rem;
    border-radius: 50%;
    border: 2px solid var(--base);
    background-size: cover;
    background-position: center;
    font-weight: 600;
    filter: drop-shadow(0px, 2px, 2px rgba(0, 0, 0.5));
    display: flex;
    align-items: center;
    justify-content: center;
    button {
      border: none;
      padding: 0;
      margin: 0;
      display: flex;
      align-items: center;
      justify-content: center;
      box-shadow: none;
      position: absolute;
      width: 2rem;
      height: 2rem;
      svg {
        width: 100%;
        height: 100%;
        fill: none;
      }
    }
    .cancel-hangup-button {
      filter: drop-shadow(var(--shadow-medium));
      bottom: 3px;
      right: 3px;
      background: red;
      border-radius: 50%;
      padding: 3px;
      border: 2px solid var(--text-color);
      transform: scaleX(-1);
    }
    .accept-button {
      filter: drop-shadow(var(--shadow-medium));
      bottom: 3px;
      left: 3px;
      background: lime;
      border-radius: 50%;
      padding: 3px;
      border: 2px solid var(--text-color);
    }
  }
}
</style>
