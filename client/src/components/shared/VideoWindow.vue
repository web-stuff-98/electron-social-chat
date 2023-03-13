<script lang="ts" setup>
import { ref } from "vue";
import { userMediaProperties } from "../../store/MediaStore";
import { userStore } from "../../store/UserStore";
defineProps<{
  userMedia?: MediaStream;
  displayMedia?: MediaStream;
  isOwner: boolean;
  uid?: string;
}>();
const hideUserMedia = ref(false);
const muteUserMedia = ref(false);
const muteDisplayMedia = ref(false);
</script>

<template>
  <div class="video-window">
    <div class="name">{{ userStore.getUser(uid as string)?.username }}</div>
    <video
      :srcObject="displayMedia || userMedia"
      :muted="
        isOwner ||
        (!displayMedia && muteUserMedia) ||
        Boolean(muteDisplayMedia && displayMedia && userMedia)
      "
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
      <button
        class="mute-button"
        @click="
          {
            if (displayMedia && userMedia) {
              muteDisplayMedia = !muteDisplayMedia;
            } else {
              muteUserMedia = !muteUserMedia;
            }
          }
        "
      >
        <v-icon
          :style="{ fill: 'white', filter:'drop-shadow(0px, 2px, 2px black)' }"
          :name="
            (displayMedia && userMedia ? muteDisplayMedia : muteUserMedia)
              ? 'bi-mic-mute-fill'
              : 'bi-mic-fill'
          "
        />
      </button>
    </div>
    <div
      :style="hideUserMedia ? { filter: 'opacity(0.5)' } : {}"
      v-if="displayMedia && userMedia"
      v-show="userMediaProperties.video"
      class="small-video-container"
    >
      <video
        v-if="!hideUserMedia"
        :srcObject="userMedia || displayMedia"
        :muted="isOwner || muteUserMedia"
        class="small-video"
        autoplay
      />
      <div class="buttons">
        <!-- Mute/unmute userMedia video button -->
        <button
          v-if="!isOwner"
          class="mute-button"
          @click="muteUserMedia = !muteUserMedia"
        >
          <v-icon :name="muteUserMedia ? 'bi-mic-mute-fill' : 'bi-mic-fill'" />
        </button>
        <!-- Hide/unhide userMedia video button -->
        <button @click="hideUserMedia = !hideUserMedia">
          <v-icon :name="hideUserMedia ? 'gi-expand' : 'io-close'" />
        </button>
      </div>
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
