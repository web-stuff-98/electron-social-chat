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
const peerUserStreamHasVideo = ref(false);
const peerDisplayStream = ref<MediaStream>();
const peerDisplayStreamHasVideo = ref(false);

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
const { userStream, displayStream, userMediaStreamID } = useChatMedia(
  negotiateConnection,
  mediaOptions
);

const peerUserMediaStreamID = ref("");

function negotiateConnection(isOnMounted?: boolean) {
  gotAnswer.value = false;
  peerUserStreamHasVideo.value = false;
  peerDisplayStreamHasVideo.value = false;
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

          um_stream_id: userMediaStreamID.value,
          um_vid: mediaOptions.value.userMedia.video,
          dm_vid: mediaOptions.value.displayMedia.video,
        })
      );
    }
  });
  peerInstance.value = peer;
}
// for recipient peer
async function makeAnswerPeer(
  signal: Peer.SignalData,
  userMediaID: string,
  showUserVid: boolean,
  showDisplayVid: boolean
) {
  const peer = initPeer();
  peer.on("signal", (signal) => {
    socketStore.send(
      JSON.stringify({
        event_type: "CALL_WEBRTC_ANSWER",
        signal: JSON.stringify(signal),

        um_stream_id: userMediaStreamID.value,
        um_vid: mediaOptions.value.userMedia.video,
        dm_vid: mediaOptions.value.displayMedia.video,
      })
    );
  });
  peerUserMediaStreamID.value = userMediaID;
  peerUserStreamHasVideo.value = showUserVid;
  peerDisplayStreamHasVideo.value = showDisplayVid;
  await nextTick(() => {
    peer.signal(signal);
  });
  peerInstance.value = peer;
}

async function signalAnswer(
  signal: Peer.SignalData,
  userMediaID: string,
  showUserVid: boolean,
  showDisplayVid: boolean
) {
  gotAnswer.value = true;
  peerUserMediaStreamID.value = userMediaID;
  peerUserStreamHasVideo.value = showUserVid;
  peerDisplayStreamHasVideo.value = showDisplayVid;
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
    makeAnswerPeer(
      JSON.parse(data.signal) as Peer.SignalData,
      data.um_stream_id,
      data.um_vid,
      data.dm_vid
    );
  }
  if (instanceOfCallWebRTCAnswerFromRecipient(data)) {
    signalAnswer(
      JSON.parse(data.signal) as Peer.SignalData,
      data.um_stream_id,
      data.um_vid,
      data.dm_vid
    );
  }
  if (instanceOfCallWebRTCRecipientRequestedReInitialization(data)) {
    console.log("Renegotiation requested");
    negotiateConnection();
  }
}

function handleStream(stream: MediaStream) {
  if (stream.id === peerUserMediaStreamID.value) {
    peerUserStream.value = stream;
  } else {
    peerDisplayStream.value = stream;
  }
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
</script>

<template>
  <div class="container">
    <div :style="{ fontSize: '0.666rem', letterSpacing: '-1px' }">
      {{ userMediaStreamID }}
      <br />
      {{ peerUserMediaStreamID }}
    </div>
    <div class="vid-chat-users">
      <!-- Current user -->
      <VidChatUser
        :userMediaStreamID="userMediaStreamID"
        :uid="authStore.user?.ID"
        :userMedia="userStream"
        :displayMedia="displayStream"
        :isOwner="true"
        :mediaOptions="mediaOptions"
        :hasDisplayMediaVideo="mediaOptions.displayMedia.video"
        :hasUserMediaVideo="mediaOptions.userMedia.video"
      />
      <!-- Other user -->
      <VidChatUser
        :userMedia="peerUserStream"
        :displayMedia="peerDisplayStream"
        :userMediaStreamID="peerUserMediaStreamID"
        :isOwner="false"
        :uid="String(otherUsersId)"
        :hasDisplayMediaVideo="peerDisplayStreamHasVideo"
        :hasUserMediaVideo="peerUserStreamHasVideo"
      />
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
  .vid-chat-users {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
    align-items: center;
    justify-content: center;
    padding: 0.6rem;
    flex-shrink: 1;
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
