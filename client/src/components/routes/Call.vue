<script lang="ts" setup>
import { authStore } from "../../store/AuthStore";
import { ref } from "vue";

const muted = ref(false);
const camera = ref(false);
const screen = ref(false);
</script>

<template>
  <div class="container">
    <div class="pfps">
      <!-- Current users pfp -->
      <div
        :style="{
          ...(authStore.user?.base64pfp
            ? { backgroundImage: `url(${authStore.user?.base64pfp})` }
            : {}),
        }"
        class="pfp"
      >
        <v-icon v-if="!authStore.user?.base64pfp" name="fa-user" />
      </div>
      <!-- Call recipients pfp -->
      <div
        :style="{
          ...(authStore.user?.base64pfp
            ? { backgroundImage: `url(${authStore.user?.base64pfp})` }
            : {}),
        }"
        class="pfp"
      >
        <v-icon v-if="!authStore.user?.base64pfp" name="fa-user" />
      </div>
    </div>
    <div class="control-buttons">
      <!-- Camera button -->
      <button @click="camera = !camera" type="button">
        <v-icon :name="camera ? 'bi-camera-video-off' : 'bi-camera-video'" />
      </button>
      <!-- Screenshare button -->
      <button @click="screen = !screen" type="button">
        <v-icon :name="screen ? 'md-stopscreenshare' : 'md-screenshare'" />
      </button>
      <!-- Mute/unmute button -->
      <button @click="muted = !muted" type="button">
        <v-icon :name="muted ? 'bi-mic-mute-fill' : 'bi-mic-fill'" />
      </button>
      <!-- Hangup button -->
      <button class="close-button" type="button">
        <v-icon name="hi-phone-missed-call" />
      </button>
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
    padding: 0.6rem;
    flex-shrink: 1;
    .pfp {
      width: 4rem;
      height: 4rem;
      border: 1px solid var(--base);
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
    border: 1px solid var(--base-light);
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
