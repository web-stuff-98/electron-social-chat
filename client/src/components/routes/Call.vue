<script lang="ts" setup>
import {
  toRef,
  onMounted,
  onBeforeUnmount,
  ref,
  watch,
  nextTick,
  computed,
} from "vue";
import { authStore } from "../../store/AuthStore";
import { useRoute, useRouter } from "vue-router";
import { userStore } from "../../store/UserStore";
import { socketStore } from "../../store/SocketStore";
import {
  closeAllMedia,
  muteMic,
  openMic,
  openCamera,
  closeCamera,
  openScreen,
  closeScreen,
  userMediaProperties,
  displayMediaActive,
  userMedia,
  displayMedia,
} from "../../store/MediaStore";
import {
  instanceOfCallLeftData,
  instanceOfCallWebRTCOfferFromInitiator,
  parseSocketEventData,
} from "../../utils/determineSocketEvent";
import VideoWindow from "../shared/VideoWindow.vue";
import Peer from "simple-peer";

/*
https://codesandbox.io/s/vinnu1simple-videochat-webrtc-0ozmn
*/

const route = useRoute();
const router = useRouter();
const otherUsersId = toRef(route.params, "id");
const initiator = computed(() => route.query.initiator !== undefined);

const PeerInstance = ref<Peer.Instance>();
const peerStream = ref<MediaStream>();

function watchForCallEvents(e: MessageEvent) {
  const data = parseSocketEventData(e);
  if (!data) return;
  if (instanceOfCallLeftData(data)) {
    PeerInstance.value?.destroy();
    router.push("/");
  }
  if (instanceOfCallWebRTCOfferFromInitiator(data)) {
    PeerInstance.value = new Peer({ initiator: false });
    PeerInstance.value.on("stream", handleStream);
    PeerInstance.value.on("signal", (signal) =>
      socketStore.send(
        JSON.stringify({
          event_type: "CALL_WEBRTC_ANSWER",
          signal,
        })
      )
    );
    PeerInstance.value.signal(JSON.parse(data.signal) as Peer.SignalData);
  }
}

function handleStream(stream: MediaStream) {
  peerStream.value = stream;
}

function handleSignal(signal: Peer.SignalData) {
  if (initiator) {
    socketStore.send(
      JSON.stringify({
        event_type: "CALL_WEBRTC_OFFER",
        signal: JSON.stringify(signal),
      })
    );
  }
}

onMounted(() => {
  socketStore.socket?.addEventListener("message", watchForCallEvents);
  openMic().then(async () => {
    if (initiator) {
      await nextTick(() => {
        // Initialize the peer
        PeerInstance.value = new Peer({ initiator: true });
        PeerInstance.value.on("stream", handleStream);
        PeerInstance.value.on("signal", handleSignal);
      });
    }
  });
});
onBeforeUnmount(() => {
  socketStore.socket?.removeEventListener("message", watchForCallEvents);
  socketStore.send(
    JSON.stringify({
      event_type: "CALL_LEAVE",
    })
  );
  closeAllMedia();
});
</script>

<template>
  <div class="container">
    <div class="pfps">
      <!-- Current users pfp / streams -->
      <div
        :style="{
          ...(authStore.user?.base64pfp
            ? { backgroundImage: `url(${authStore.user?.base64pfp})` }
            : {}),
        }"
        class="pfp"
        v-if="!userMediaProperties.video && !displayMediaActive"
      >
        <v-icon v-if="!authStore.user?.base64pfp" name="fa-user" />
      </div>
      <VideoWindow
        v-else
        :uid="authStore.user?.ID"
        :userMedia="userMedia"
        :displayMedia="displayMedia"
        :isOwner="true"
      />
      <!-- Other users pfp / streams -->
      <div
        :style="{
          ...(userStore.getUser(otherUsersId as string)?.base64pfp
            ? { backgroundImage: `url(${userStore.getUser(otherUsersId as string)?.base64pfp})` }
            : {}),
        }"
        class="pfp"
        v-if="!peerStream"
      >
        <v-icon
          v-if="!userStore.getUser(otherUsersId as string)?.base64pfp"
          name="fa-user"
        />
      </div>
      <VideoWindow
        v-else
        :uid="authStore.user?.ID"
        :userMedia="peerStream"
        :isOwner="false"
      />
    </div>
    <div class="control-buttons">
      <!-- Camera button -->
      <button
        @click="
          {
            if (!userMediaProperties.video) {
              openCamera();
            } else {
              closeCamera();
            }
          }
        "
        type="button"
      >
        <v-icon
          :name="
            userMediaProperties.video
              ? 'bi-camera-video-off'
              : 'bi-camera-video'
          "
        />
      </button>
      <!-- Screenshare button -->
      <button
        @click="
          {
            if (!displayMediaActive) {
              openScreen();
            } else {
              closeScreen();
            }
          }
        "
        type="button"
      >
        <v-icon
          :name="displayMediaActive ? 'md-stopscreenshare' : 'md-screenshare'"
        />
      </button>
      <!-- Mute/unmute button -->
      <button
        @click="
          {
            if (userMediaProperties.audio) {
              muteMic();
            } else {
              openMic();
            }
          }
        "
        type="button"
      >
        <v-icon
          :name="
            !userMediaProperties.audio ? 'bi-mic-mute-fill' : 'bi-mic-fill'
          "
        />
      </button>
      <!-- Hangup button -->
      <router-link to="/">
        <button class="close-button" type="button">
          <v-icon name="hi-phone-missed-call" />
        </button>
      </router-link>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.container {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  .pfps {
    display: flex;
    gap: 1vw;
    flex-wrap: wrap;
    align-items: center;
    justify-content: center;
    padding: 0.6rem;
    flex-shrink: 1;
    .pfp {
      width: 4rem;
      height: 4rem;
      border: 2px solid var(--base);
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      box-shadow: var(--shadow);
      gap: var(--padding);
      background-size: cover;
      background-position: center;
      svg {
        width: 2rem;
        height: 2rem;
      }
    }
  }
  .control-buttons {
    display: flex;
    gap: 0.5rem;
    padding: 0.25rem;
    border: 2px solid var(--base-light);
    border-radius: 5pc;
    box-shadow: var(--shadow);
    button {
      border: 2px solid var(--base);
      border-radius: 50%;
      padding: 0;
      margin: 0;
      width: 2.5rem;
      height: 2.5rem;
      box-shadow: var(--shadow);
      display: flex;
      align-items: center;
      justify-content: center;
      svg {
        width: 70%;
        height: 70%;
      }
    }
    .close-button {
      background: red;
      border: 2px solid var(--text-color);
      svg {
        fill: none;
        margin-right: 0.1rem;
        margin-top: 0.1rem;
      }
    }
  }
}
</style>
