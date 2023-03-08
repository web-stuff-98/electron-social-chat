<script lang="ts" setup>
import { authStore } from "../../store/AuthStore";
</script>

<template>
  <div v-if="authStore.user" class="container">
    <h1>Welcome</h1>
    <h2>To get started click on the menu icon in the bottom left</h2>
  </div>
</template>

<style lang="scss" scoped>
.container {
  display: flex;
  flex-direction: column;
  padding: calc(2 * var(--padding));
  gap: var(--padding-medium);
  width: 100%;
  height: 100%;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  pointer-events: none;
  background: var(--background-radial);
  filter: drop-shadow(0px 3px 4px black);
  h1,
  h2 {
    margin: 0;
  }
  h2 {
    margin-top: var(--padding);
  }
}
</style>
