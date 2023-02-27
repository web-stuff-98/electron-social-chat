<script setup lang="ts">
import { ref } from "vue";

import Profile from "./sections/Profile.vue";
import CreateRoom from "./sections/CreateRoom.vue";
import ExploreRooms from "./sections/ExploreRooms.vue";
import DirectMessages from "./sections/DirectMessages.vue";

enum EAsideSection {
  "PROFILE" = "Profile",
  "CREATE_ROOM" = "Create room",
  "EXPLORE_ROOMS" = "Explore rooms",
  "DIRECT_MESSAGES" = "Direct messages",
}

const show = ref(false);
const section = ref<EAsideSection>(EAsideSection.PROFILE);
</script>

<template>
  <aside v-show="show">
    <div class="buttons">
      <button @click="section = EAsideSection.DIRECT_MESSAGES">
        Direct messages
        <span/>
      </button>
      <button @click="section = EAsideSection.EXPLORE_ROOMS">
        Explore rooms
      </button>
      <button @click="section = EAsideSection.CREATE_ROOM">Create room</button>
      <button @click="section = EAsideSection.PROFILE">Profile</button>
    </div>
    <div class="container">
      <div class="content">
        <DirectMessages v-if="section === EAsideSection.DIRECT_MESSAGES" />
        <ExploreRooms v-if="section === EAsideSection.EXPLORE_ROOMS" />
        <CreateRoom v-if="section === EAsideSection.CREATE_ROOM" />
        <Profile v-if="section === EAsideSection.PROFILE" />
      </div>
      <button @click="show = false" class="close-button">
        <v-icon name="io-close" />
      </button>
    </div>
  </aside>
  <button v-if="!show" @click="show = true" class="aside-menu-button">
    <v-icon class="aside-menu-icon" name="co-menu" />
  </button>
</template>

<style lang="scss" scoped>
aside {
  box-sizing: border-box;
  width: 10rem;
  height: 100%;
  max-height: calc(100% - var(--header-height));
  border-right: 1px solid var(--base-light);
  display: flex;
  flex-direction: column;
  padding: 0;
  background: var(--foreground);
  box-shadow: 0px 0px 3px black;
  .buttons {
    button {
      border: none;
      border-radius: 0;
      width: 100%;
      padding: var(--padding-medium);
      text-align: left;
      border-bottom: 1px solid var(--base-light);
      box-shadow: none;
      font-size: 1rem;
      position: relative;
      span {
        width: 2px;
        height: 2px;
        position: absolute;
        top: 3px;
        right: 3px;
      }
    }
  }
  .close-button {
    border: 1px solid white;
    padding: 0;
    margin: 0;
    filter: opacity(0.666);
    width: fit-content;
    box-shadow: var(--shadow-medium);
    transition: filter 100ms ease;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--border-radius-medium);
    svg {
      width: 1rem;
      height: 1rem;
    }
  }
  .close-button:hover {
    outline: 1px solid white;
    filter: opacity(1);
  }
}
.aside-menu-button {
  border: none;
  background: none;
  padding: 0;
  position: fixed;
  bottom: 3px;
  left: 3px;
}
.aside-menu-icon {
  width: 1.666rem;
  height: 1.666rem;
  z-index: 99;
  border: 2px solid white;
  border-radius: var(--border-radius-medium);
}
.container {
  padding: var(--padding-medium);
  align-items: flex-end;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  padding: var(--padding-medium);
  box-sizing: border-box;
  .content {
    flex-grow: 1;
    width: 100%;
    height: 100%;
  }
}
</style>
