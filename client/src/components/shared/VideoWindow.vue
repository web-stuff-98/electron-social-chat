<script lang="ts" setup>
import { ref, toRefs, watch } from "vue";
import { userStore } from "../../store/UserStore";
const props = defineProps<{
  media: MediaStream | undefined;
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
</script>

<template>
  <div
    v-show="trackIds.userMediaVideo || trackIds.displayMediaVideo"
    class="video-window"
  >
    <div class="name">
      {{ userStore.getUser(uid as string)?.username }}
    </div>
    <video
      v-show="trackIds.userMediaVideo"
      :srcObject="media"
      :muted="isOwner"
      class="main-video"
      autoplay
    />
    <div
      v-if="(!isOwner && trackIds.userMediaVideo) || trackIds.displayMediaVideo"
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
