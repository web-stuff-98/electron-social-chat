<script lang="ts" setup>
import { authStore } from "../../store/AuthStore";
</script>

<template>
  <div v-if="authStore.user" class="container">
    <h1>Welcome</h1>
    <span />
    <h2>To get started click on the menu icon in the bottom left</h2>
    <p>
      <b>Hint:</b> you can click on other users to open an action menu dropdown
    </p>
  </div>
</template>

<style lang="scss" scoped>
.container {
  display: flex;
  flex-direction: column;
  padding: calc(2 * var(--padding));
  gap: var(--padding-medium);
  width: 100%;
  max-width: 50rem;
  height: 100%;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  pointer-events: none;
  filter: drop-shadow(var(--shadow-medium));
  span {
    margin: 0;
    padding: 0;
    border-bottom: 2px solid var(--text-color);
    width: 100%;
    mask-image: linear-gradient(
      to right,
      transparent 0%,
      black 50%,
      transparent 100%
    );
    -webkit-mask-image: linear-gradient(
      to right,
      transparent 0%,
      black 50%,
      transparent 100%
    );
  }
  h1,
  h2,
  p {
    margin: 0;
  }
}
</style>
