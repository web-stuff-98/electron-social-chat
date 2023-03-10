<script setup lang="ts">
import { ref } from "vue";

import Profile from "./sections/Profile.vue";
import CreateRoom from "./sections/CreateRoom.vue";
import ExploreRooms from "./sections/ExploreRooms.vue";
import DirectMessages from "./sections/DirectMessages.vue";
import FindUser from "./sections/FindUser.vue";
import { authStore } from "../../../store/AuthStore";
import { showCreateRoom } from "../../../store/CreateEditRoomStore";

enum EAsideSection {
  "PROFILE" = "Profile",
  "EXPLORE_ROOMS" = "Explore rooms",
  "YOUR_ROOMS" = "Your rooms",
  "DIRECT_MESSAGES" = "Direct messages",
  "FIND_USER" = "Find user",
}

const section = ref<EAsideSection>(EAsideSection.PROFILE);
defineProps<{ showAside: boolean; toggleShowAside: Function }>();

function logout() {
  authStore.logout();
}
</script>

<template>
  <aside
    :style="{
      transform: `translateX(${showAside ? '0%' : '-100%'})`,
      ...(showAside ? {} : { transition: 'none' }),
    }"
  >
    <div class="buttons">
      <button @click="section = EAsideSection.DIRECT_MESSAGES">
        Direct messages
        <span v-if="section === EAsideSection.DIRECT_MESSAGES" />
      </button>
      <button @click="section = EAsideSection.EXPLORE_ROOMS">
        Explore rooms
        <span v-if="section === EAsideSection.EXPLORE_ROOMS" />
      </button>
      <button @click="section = EAsideSection.YOUR_ROOMS">
        Your rooms
        <span v-if="section === EAsideSection.YOUR_ROOMS" />
      </button>
      <button @click="showCreateRoom = true">
        Create room
        <span v-if="showCreateRoom" />
      </button>
      <button @click="section = EAsideSection.PROFILE">
        Profile
        <span v-if="section === EAsideSection.PROFILE" />
      </button>
      <button @click="section = EAsideSection.FIND_USER">
        Find user
        <span v-if="section === EAsideSection.FIND_USER" />
      </button>
      <button @click="logout">Log out</button>
    </div>
    <div class="container">
      <div class="content">
        <DirectMessages v-if="section === EAsideSection.DIRECT_MESSAGES" />
        <ExploreRooms
          :own="false"
          v-if="section === EAsideSection.EXPLORE_ROOMS"
        />
        <Profile v-if="section === EAsideSection.PROFILE" />
        <FindUser v-if="section === EAsideSection.FIND_USER" />
        <ExploreRooms :own="true" v-if="section === EAsideSection.YOUR_ROOMS" />
      </div>
      <button @click="() => toggleShowAside()" class="close-button">
        <v-icon name="io-close" />
      </button>
    </div>
  </aside>
  <button
    v-if="!showAside"
    @click="() => toggleShowAside()"
    class="aside-menu-button"
  >
    <v-icon class="aside-menu-icon" name="co-menu" />
  </button>
</template>

<style lang="scss" scoped>
aside {
  box-sizing: border-box;
  min-width: var(--aside-width);
  width: var(--aside-width);
  height: 100vh;
  border-right: 2px solid var(--base-light);
  display: flex;
  flex-direction: column;
  background: var(--foreground);
  transition: transform 50ms linear;
  box-shadow: 2px 0px 2px rgba(0, 0, 0, 0.02);
  padding-right: 1px;
  padding-top: 2px;
  position: fixed;
  left: 0;
  .buttons {
    display: flex;
    flex-direction: column;
    padding: var(--padding-medium);
    padding-bottom: 0;
    padding-left: calc(var(--padding-medium) + 2px);
    padding-top: calc(var(--padding-medium) + 1px);
    gap: var(--padding-medium);
    button {
      width: 100%;
      border: 2px solid var(--base-light);
      padding: var(--padding-medium);
      text-align: left;
      font-size: 1rem;
      position: relative;
      span {
        width: 3px;
        height: 3px;
        position: absolute;
        top: 3px;
        border-radius: 50%;
        right: 3px;
        background: var(--text-color);
        box-shadow: 0px 0px 3px var(--text-color), 0px 0px 8px var(--text-color);
      }
    }
  }
  .close-button {
    border: 1px solid var(--text-color);
    padding: 0;
    margin: var(--padding-medium);
    margin-bottom: calc(var(--padding-medium) + 1px);
    filter: opacity(0.666);
    width: fit-content;
    box-shadow: var(--shadow-medium);
    transition: filter 100ms ease;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--border-radius-medium);
    background: red;
    svg {
      width: 1rem;
      height: 1rem;
    }
  }
  .close-button:hover {
    outline: 1px solid var(--text-color);
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
  z-index: 100;
  filter: opacity(0.5);
}
.aside-menu-button:hover {
  filter: opacity(1);
}
.aside-menu-icon {
  width: 1.25rem;
  height: 1.25rem;
  z-index: 99;
  border: 2px solid var(--text-color);
  border-radius: var(--border-radius-medium);
}
.container {
  padding: calc(var(--padding-medium) - 2px);
  padding-left: calc(var(--padding-medium) - 1px);
  align-items: flex-end;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  box-sizing: border-box;
  height: 100%;
  .content {
    box-sizing: border-box;
    flex-grow: 1;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
</style>
