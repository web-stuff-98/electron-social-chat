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
import VideoWindow from "../shared/VideoWindow.vue";
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
const peerStream = ref<MediaStream>();
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
const { stream, trackIds } = useChatMedia(negotiateConnection, mediaOptions);

type TrackIDs = {
  um_snd_track_id: string;
  um_vid_track_id: string;
  dm_snd_track_id: string;
  dm_vid_track_id: string;
};

const peerTrackIds = ref({
  userMediaAudio: "",
  userMediaVideo: "",
  displayMediaAudio: "",
  displayMediaVideo: "",
});

function negotiateConnection(isOnMounted?: boolean) {
  gotAnswer.value = false;
  if (initiator.value) {
    console.log("Is negotiating as initiator");
    if (peerInstance.value) {
      peerInstance.value.destroy();
    }
    peerStream.value = undefined;
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
    stream: stream.value,
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

          um_snd_track_id: trackIds.userMediaAudio,
          um_vid_track_id: trackIds.userMediaVideo,
          dm_snd_track_id: trackIds.displayMediaAudio,
          dm_vid_track_id: trackIds.displayMediaVideo,
        })
      );
    }
  });
  peerInstance.value = peer;
}
// for recipient peer
async function makeAnswerPeer(signal: Peer.SignalData, pTrackIds: TrackIDs) {
  const peer = initPeer();
  peer.on("signal", (signal) => {
    socketStore.send(
      JSON.stringify({
        event_type: "CALL_WEBRTC_ANSWER",
        signal: JSON.stringify(signal),

        um_snd_track_id: trackIds.userMediaAudio,
        um_vid_track_id: trackIds.userMediaVideo,
        dm_snd_track_id: trackIds.displayMediaAudio,
        dm_vid_track_id: trackIds.displayMediaVideo,
      })
    );
  });
  peerTrackIds.value = {
    userMediaAudio: pTrackIds.um_snd_track_id,
    userMediaVideo: pTrackIds.um_vid_track_id,
    displayMediaAudio: pTrackIds.dm_snd_track_id,
    displayMediaVideo: pTrackIds.dm_vid_track_id,
  };
  await nextTick(() => {
    peer.signal(signal);
  });
  peerInstance.value = peer;
}

async function signalAnswer(signal: Peer.SignalData, pTrackIds: TrackIDs) {
  gotAnswer.value = true;
  peerTrackIds.value = {
    userMediaAudio: pTrackIds.um_snd_track_id,
    userMediaVideo: pTrackIds.um_vid_track_id,
    displayMediaAudio: pTrackIds.dm_snd_track_id,
    displayMediaVideo: pTrackIds.dm_vid_track_id,
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
      um_snd_track_id: data.um_snd_track_id,
      um_vid_track_id: data.um_vid_track_id,
      dm_snd_track_id: data.dm_snd_track_id,
      dm_vid_track_id: data.dm_vid_track_id,
    });
  }
  if (instanceOfCallWebRTCAnswerFromRecipient(data)) {
    signalAnswer(JSON.parse(data.signal) as Peer.SignalData, {
      um_snd_track_id: data.um_snd_track_id,
      um_vid_track_id: data.um_vid_track_id,
      dm_snd_track_id: data.dm_snd_track_id,
      dm_vid_track_id: data.dm_vid_track_id,
    });
  }
  if (instanceOfCallWebRTCRecipientRequestedReInitialization(data)) {
    console.log("Renegotiation requested");
    negotiateConnection();
  }
}

function handleStream(stream: MediaStream) {
  peerStream.value = stream;
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
  const tracks = stream.value?.getVideoTracks();
  tracks?.forEach((track) => {
    if (track.enabled) return true;
  });
  return false;
});

const showPeerVideoWindow = computed(() => {
  const tracks = peerStream.value?.getVideoTracks();
  tracks?.forEach((track) => {
    if (!track.muted) return true;
  });
  return false;
});
</script>

<template>
  <div class="container">
    <div :style="{ fontSize: '0.7rem' }">
      {{ trackIds }}
      <br />
      {{ peerTrackIds }}
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
      <VideoWindow
        :trackIds="trackIds"
        :uid="authStore.user?.ID"
        :media="stream"
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
      <VideoWindow
        v-else
        :trackIds="peerTrackIds"
        :uid="otherUsersId as string"
        :media="peerStream"
        :isOwner="false"
      />
      <video :style="{ maxWidth: '4rem' }" :srcObject="peerStream" autoplay />
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
