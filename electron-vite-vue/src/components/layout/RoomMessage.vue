<script lang="ts" setup>
import { IRoomMessage } from "../../interfaces/GeneralInterfaces";
import User from "../shared/User.vue";
defineProps<{ msg: IRoomMessage; reverse?: boolean }>();
</script>

<template>
  <div
    :style="reverse ? { textAlign: 'right', flexDirection: 'row-reverse' } : {}"
    class="message"
  >
    <div class="user-container">
      <User :dateTime="new Date(msg.created_at)" :reverse="reverse" :uid="msg.author" />
    </div>
    {{ msg.content }}
  </div>
</template>

<style lang="scss" scoped>
.message {
  text-align: left;
  font-size: 0.8rem;
  display: flex;
  align-items: center;
  .user-container {
    padding: var(--padding);
  }
}
</style>
