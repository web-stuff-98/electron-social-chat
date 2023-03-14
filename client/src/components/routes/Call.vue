<script lang="ts" setup>
import {
  toRef,
  onMounted,
  onBeforeUnmount,
  ref,
  nextTick,
  computed,
} from "vue";
import { authStore } from "../../store/AuthStore";
import { useRoute, useRouter } from "vue-router";
import { userStore } from "../../store/UserStore";
import { socketStore } from "../../store/SocketStore";
import {
  instanceOfCallLeftData,
  instanceOfCallWebRTCOfferFromInitiator,
  instanceOfCallWebRTCAnswerFromRecipient,
  instanceOfCallWebRTCRecipientRequestedReInitialization,
  parseSocketEventData,
} from "../../utils/determineSocketEvent";
import VidChatUser from "../shared/VidChatUser.vue";
import Peer from "simple-peer";
import { useChatMedia } from "../../composables/useChatMedia";

/*
basically copied this for the events
https://codesandbox.io/s/vinnu1simple-videochat-webrtc-0ozmn
*/

const route = useRoute();
const router = useRouter();
const otherUsersId = toRef(route.params, "id");
const initiator = computed(() => route.query.initiator !== undefined);

const peerInstance = ref<Peer.Instance>();
const peerUserStream = ref<MediaStream>();
const peerDisplayStream = ref<MediaStream>();
const gotAnswer = ref(false);

const mediaOptions = ref({
  userMedia: {
    audio: true,
    video: false,
  },
  displayMedia: {
    audio: false,
    video: false,
  },
});
const { userStream, displayStream, streamIds } = useChatMedia(
  negotiateConnection,
  mediaOptions
);

type StreamIDs = {
  um_stream_id: string;
  dm_stream_id: string;
};

const peerStreamIds = ref({
  userMedia: "",
  displayMedia: "",
});

function negotiateConnection(isOnMounted?: boolean) {
  gotAnswer.value = false;
  if (initiator.value) {
    console.log("Is negotiating as initiator");
    if (peerInstance.value) {
      peerInstance.value.destroy();
    }
    peerUserStream.value = undefined;
    peerDisplayStream.value = undefined;
    makePeer();
  } else if (!isOnMounted) {
    console.log("Is requesting reinitialization");
    requestReInitialization();
  }
}

function requestReInitialization() {
  socketStore.send(
    JSON.stringify({
      event_type: "CALL_WEBRTC_RECIPIENT_REQUEST_REINITIALIZATION",
    })
  );
}

function initPeer() {
  const peer = new Peer({
    initiator: initiator.value,
    trickle: false,
    streams: [
      ...(userStream.value ? [userStream.value] : []),
      ...(displayStream.value ? [displayStream.value] : []),
    ],
    iceCompleteTimeout: 2000, // 5 seconds is too long
  });
  peer.on("stream", handleStream);
  return peer;
}

// for initializer peer
function makePeer() {
  gotAnswer.value = false;
  const peer = initPeer();
  peer.on("signal", (signal) => {
    if (!gotAnswer.value) {
      socketStore.send(
        JSON.stringify({
          event_type: "CALL_WEBRTC_OFFER",
          signal: JSON.stringify(signal),

          um_stream_id: streamIds.userMedia,
          dm_stream_id: streamIds.displayMedia,
        })
      );
    }
  });
  peerInstance.value = peer;
}
// for recipient peer
async function makeAnswerPeer(signal: Peer.SignalData, pStreamIds: StreamIDs) {
  const peer = initPeer();
  peer.on("signal", (signal) => {
    socketStore.send(
      JSON.stringify({
        event_type: "CALL_WEBRTC_ANSWER",
        signal: JSON.stringify(signal),

        um_stream_id: streamIds.userMedia,
        dm_stream_id: streamIds.displayMedia,
      })
    );
  });
  peerStreamIds.value = {
    userMedia: pStreamIds.um_stream_id,
    displayMedia: pStreamIds.dm_stream_id,
  };
  await nextTick(() => {
    peer.signal(signal);
  });
  peerInstance.value = peer;
}

async function signalAnswer(signal: Peer.SignalData, pStreamIds: StreamIDs) {
  gotAnswer.value = true;
  peerStreamIds.value = {
    userMedia: pStreamIds.um_stream_id,
    displayMedia: pStreamIds.dm_stream_id,
  };
  await nextTick(() => {
    peerInstance.value?.signal(signal);
  });
}

function watchForCallEvents(e: MessageEvent) {
  const data = parseSocketEventData(e);
  if (!data) return;
  if (instanceOfCallLeftData(data)) {
    peerInstance.value?.destroy();
    router.push("/");
  }
  if (instanceOfCallWebRTCOfferFromInitiator(data)) {
    makeAnswerPeer(JSON.parse(data.signal) as Peer.SignalData, {
      um_stream_id: data.um_stream_id,
      dm_stream_id: data.dm_stream_id,
    });
  }
  if (instanceOfCallWebRTCAnswerFromRecipient(data)) {
    signalAnswer(JSON.parse(data.signal) as Peer.SignalData, {
      um_stream_id: data.um_stream_id,
      dm_stream_id: data.dm_stream_id,
    });
  }
  if (instanceOfCallWebRTCRecipientRequestedReInitialization(data)) {
    console.log("Renegotiation requested");
    negotiateConnection();
  }
}

function handleStream(stream: MediaStream) {
  if (stream.id === peerStreamIds.value.userMedia)
    peerUserStream.value = stream;
  if (stream.id === peerStreamIds.value.displayMedia)
    peerDisplayStream.value = stream;
}

onMounted(() => {
  socketStore.socket?.addEventListener("message", watchForCallEvents);
  userStore.cacheUserData(otherUsersId.value as string, true);
});
onBeforeUnmount(() => {
  socketStore.socket?.removeEventListener("message", watchForCallEvents);
  socketStore.send(
    JSON.stringify({
      event_type: "CALL_LEAVE",
    })
  );
});

const showVideoWindow = computed(() => {
  return streamIds.displayMedia || streamIds.userMedia;
});

const showPeerVideoWindow = computed(() => {
  return peerStreamIds.value.displayMedia || peerStreamIds.value.userMedia;
});
</script>

<template>
  <div class="container">
    <div :style="{ fontSize: '0.7rem' }">
      {{ streamIds }}
      <br />
      {{ peerStreamIds }}
    </div>
    <div class="pfps">
      <!-- Current users pfp / streams -->
      <div
        :style="{
          ...(authStore.user?.base64pfp
            ? { backgroundImage: `url(${authStore.user?.base64pfp})` }
            : {}),
        }"
        class="pfp"
        v-if="!showVideoWindow"
      >
        <v-icon v-if="!authStore.user?.base64pfp" name="fa-user" />
      </div>
      <VidChatUser
        :streamIds="streamIds"
        :uid="authStore.user?.ID"
        :userMedia="userStream"
        :displayMedia="displayStream"
        :isOwner="true"
        :mediaOptions="mediaOptions"
      />
      <!-- Other users pfp / streams -->
      <div
        :style="{
          ...(userStore.getUser(otherUsersId as string)?.base64pfp
            ? { backgroundImage: `url(${userStore.getUser(otherUsersId as string)?.base64pfp})` }
            : {}),
        }"
        class="pfp"
        v-if="!showPeerVideoWindow"
      >
        <v-icon
          v-if="!userStore.getUser(otherUsersId as string)?.base64pfp"
          name="fa-user"
        />
      </div>
      <VidChatUser
        :userMedia="peerUserStream"
        :displayMedia="peerDisplayStream"
        :streamIds="peerStreamIds"
        :isOwner="false"
        :uid="String(otherUsersId)"
      />
      <!-- <video :style="{ maxWidth: '4rem' }" :srcObject="peerStream" autoplay /> -->
    </div>
    <div class="control-buttons">
      <!-- Camera button -->
      <button
        @click="mediaOptions.userMedia.video = !mediaOptions.userMedia.video"
        type="button"
      >
        <v-icon
          :name="
            mediaOptions.userMedia.video
              ? 'bi-camera-video-off'
              : 'bi-camera-video'
          "
        />
      </button>
      <!-- Screenshare button -->
      <button
        @click="
          {
            mediaOptions.displayMedia.video = !mediaOptions.displayMedia.video;
          }
        "
        type="button"
      >
        <v-icon
          :name="
            mediaOptions.displayMedia.video
              ? 'md-stopscreenshare'
              : 'md-screenshare'
          "
        />
      </button>
      <!-- Mute/unmute button -->
      <button
        @click="
          {
            mediaOptions.userMedia.audio = !mediaOptions.userMedia.audio;
          }
        "
        type="button"
      >
        <v-icon
          :name="
            mediaOptions.userMedia.audio ? 'bi-mic-mute-fill' : 'bi-mic-fill'
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
      width: 6rem;
      height: 6rem;
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
