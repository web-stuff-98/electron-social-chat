<script lang="ts" setup>
function toggleDarkMode() {
    document.body.classList.toggle("dark-mode")
}
</script>

<template>
  <button class="container">
    <v-icon @click="toggleDarkMode" name="md-darkmode" />
  </button>
</template>

<style lang="scss" scoped>
.container {
  position: fixed;
  top: 0;
  right: 0;
  background: none;
  box-shadow: none;
  border: none;
  padding: 0;
  display: flex;
  z-index: 100;
  filter: drop-shadow(var(--shadow-medium));
  svg {
    width: 2rem;
    height: 2rem;
    fill: var(--text-color);
    color: var(--text-color);
  }
}
</style>
