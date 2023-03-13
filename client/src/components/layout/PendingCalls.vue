<script lang="ts" setup>
import { pendingCallsStore } from "../../store/PendingCallsStore";
</script>

<template>
  <div class="container">
    {{ pendingCallsStore }}
  </div>
</template>

<style lang="scss" scoped>
.container {
  position: fixed;
  bottom: 0;
  right: 0;
}
</style>
