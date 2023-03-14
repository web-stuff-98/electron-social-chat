<script lang="ts" setup>
import { computed, ref, toRefs, watch } from "vue";
import useUser from "../../composables/useUser";
import { userStore } from "../../store/UserStore";
const props = defineProps<{
  userMedia: MediaStream | undefined;
  displayMedia: MediaStream | undefined;
  isOwner: boolean;
  uid?: string;
  userMediaStreamID:string;
  // This should only be used for the current users container, not other peers.
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
  hasDisplayMediaVideo: boolean;
  hasUserMediaVideo: boolean;
}>();
const { userMedia, displayMedia, userMediaStreamID, uid, isOwner } = toRefs(props);

const user = useUser(uid?.value as string);
</script>

<template>
  <div v-show="userMediaStreamID" class="container">
    <!-- Pfp container - For when there are no video streams present -->
    <div
      :style="{
        ...(user?.base64pfp
          ? { backgroundImage: `url(${user?.base64pfp})` }
          : {}),
      }"
      class="pfp"
      v-if="!hasDisplayMediaVideo && !hasUserMediaVideo"
    >
      <v-icon v-if="!user?.base64pfp" name="fa-user" />
    </div>
    <!-- Video container - For when there is either or both video streams present -->
    <div
      v-show="hasDisplayMediaVideo || hasUserMediaVideo"
      class="vid-container"
    >
      <div class="name">
        {{ userStore.getUser(uid as string)?.username }}
      </div>
      <video
        v-show="hasDisplayMediaVideo || hasUserMediaVideo"
        :srcObject="displayMedia || userMedia"
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
            :style="{
              fill: 'white',
              filter: 'drop-shadow(0px, 2px, 2px black)',
            }"
            name="bi-mic-fill"
          />
        </button>
      </div>
      <!-- Smaller video, for when display media is present -->
      <div class="small-video-container">
        <video
          v-show="hasDisplayMediaVideo && hasUserMediaVideo"
          :srcObject="userMedia"
          :muted="isOwner"
          autoplay
        />
        <div class="buttons">
          <!-- Mute/unmute button -->
          <button class="mute-button">
            <v-icon
              :style="{
                fill: 'white',
                filter: 'drop-shadow(0px, 2px, 2px black)',
              }"
              name="bi-mic-fill"
            />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.container {
  position: relative;
  display: flex;
  flex-direction: column;
  text-align: left;
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
  .vid-container {
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
}
</style>
