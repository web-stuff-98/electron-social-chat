<script lang="ts" setup>
import { authStore } from "../../store/AuthStore";
</script>

<template>
  <div v-if="authStore.user" class="container">
    <h1>Welcome</h1>
    <hr />
    <h2>To get started click on the menu icon in the bottom left</h2>
  </div>
</template>

<style lang="scss" scoped>
.container {
  display: flex;
  flex-direction: column;
  padding: var(--padding);
  gap: var(--padding-medium);
  width: 100%;
  height: 100%;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  pointer-events: none;
  background: var(--foreground);
  filter: drop-shadow(0px 3px 4px black);
  h1,
  h2 {
    margin: 0;
  }
  hr {
    margin: 0;
    border: none;
    border-bottom: 2px solid var(--base);
    --webkit-mask-image: linear-gradient(
      to right,
      transparent 0%,
      black 50%,
      transparent 100%
    );
  }
  h2 {
    margin-top: var(--padding);
  }
}
</style>
