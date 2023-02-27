<script lang="ts">
import { ref } from "vue";
import DirectMessages from "./sections/DirectMessages.vue";
import Profile from "./sections/Profile.vue";

export default {
  setup() {
    const show = ref(false);
    return {
      show,
    };
  },
  components: { Profile, DirectMessages },
};
</script>

<template>
  <aside v-show="show">
    <div class="buttons">
      <button>Direct messages</button>
      <button>Explore rooms</button>
      <button>Create room</button>
      <button>Profile</button>
    </div>
    <div class="container">
      <div class="content">
        <DirectMessages />
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

<style scoped>
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
}
aside .buttons button {
  border: none;
  border-radius: 0;
  width: 100%;
  padding: var(--padding-medium);
  text-align: left;
  border-bottom: 1px solid var(--base-light);
  box-shadow: none;
  font-size: 1rem;
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
}
.container .content {
  flex-grow: 1;
  width: 100%;
  height: 100%;
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
}
.close-button:hover {
  outline: 1px solid white;
  filter: opacity(1);
}
.close-button svg {
  width: 1rem;
  height: 1rem;
}
</style>
