<script lang="ts" setup>
import { computed, toRefs } from "vue";
import { userStore } from "../../store/UserStore";
const props = defineProps<{
  media?: MediaStream;
  isOwner: boolean;
  uid?: string;
  trackIds: {
    userMediaVideo: string;
    userMediaAudio: string;
    displayMediaVideo: string;
    displayMediaAudio: string;
  };
  // This should only be used for the current users window, not other peers.
  // Used to enable/disable microphone and camera access
  mediaOptions?: {
    userMedia: {
      audio: boolean;
      video: boolean;
    };
    displayMedia: {
      audio: boolean;
      video: boolean;
    };
  };
}>();
const { media, trackIds } = toRefs(props);

const userMedia = computed(() => {
  const videoTrack = media?.value?.getTrackById(trackIds.value.userMediaVideo);
  const audioTrack = media?.value?.getTrackById(trackIds.value.userMediaAudio);
  const stream = new MediaStream();
  if (videoTrack) stream.addTrack(videoTrack);
  if (audioTrack) stream.addTrack(audioTrack);
  return videoTrack || audioTrack ? stream : undefined;
});

const displayMedia = computed(() => {
  const videoTrack = media?.value?.getTrackById(
    trackIds.value.displayMediaVideo
  );
  const audioTrack = media?.value?.getTrackById(
    trackIds.value.displayMediaAudio
  );
  const stream = new MediaStream();
  if (videoTrack) stream.addTrack(videoTrack);
  if (audioTrack) stream.addTrack(audioTrack);
  return videoTrack || audioTrack ? stream : undefined;
});

const showUserMediaVideo = computed(() => {
  const vidTracks = userMedia.value?.getVideoTracks();
  return vidTracks ? vidTracks[0].enabled : false;
});

const showDisplayMediaVideo = computed(() => {
  const vidTracks = displayMedia.value?.getVideoTracks();
  return vidTracks ? vidTracks[0].enabled : false;
});
</script>

<template>
  <div class="video-window">
    <div class="name">{{ userStore.getUser(uid as string)?.username }}</div>
    <video
      v-show="showUserMediaVideo"
      :srcObject="userMedia"
      :muted="isOwner"
      class="main-video"
      autoplay
    />
    <div
      v-if="!isOwner"
      :style="{
        width: 'fit-content',
        position: 'absolute',
        bottom: '1rem',
        left: '1rem',
      }"
      class="buttons"
    >
      <!-- Mute/unmute button -->
      <button class="mute-button">
        <v-icon
          :style="{ fill: 'white', filter: 'drop-shadow(0px, 2px, 2px black)' }"
          name="bi-mic-fill"
        />
      </button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.video-window {
  position: relative;
  padding: var(--padding);
  display: flex;
  flex-direction: column;
  text-align: left;
  .buttons {
    display: flex;
    justify-content: flex-end;
    height: 1.5rem;
    button svg {
      width: 100%;
      height: 100%;
    }
    .mute-button {
      svg {
        width: 70%;
        height: 70%;
      }
    }
    button {
      height: 1.5rem;
      width: 1.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 0;
      margin: 0;
      border: none;
      box-shadow: none;
      background: none;
    }
    video {
      width: 100%;
    }
  }
  .name {
    padding: var(--padding-medium);
    position: absolute;
    top: var(--padding);
    left: var(--padding);
    padding: var(--padding-medium);
    font-weight: 600;
    text-shadow: 0px 2px 2px black;
    color: white;
  }
  .main-video,
  .small-video-container {
    border: 1px solid var(--base-light);
    height: auto;
    box-shadow: var(--shadow);
    border-radius: var(--border-radius-medium);
  }
  .main-video {
    width: 45vw;
    max-width: min(30rem, 40vh);
  }
  .small-video-container {
    position: absolute;
    bottom: 0;
    right: 0;
    width: 30%;
    height: auto;
    background: var(--foreground);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
}
</style>
