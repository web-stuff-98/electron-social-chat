<script setup lang="ts">
import { IResMsg } from "../../interfaces/GeneralInterfaces";
defineProps<{ resMsg: IResMsg }>();
</script>

<template>
  <div v-show="resMsg.pen || resMsg.err || resMsg.msg" class="res-msg">
    <v-icon v-if="resMsg.err" class="res-msg-err-icon" name="md-error-round" />
    <p :style="resMsg.err ? { color: 'red' } : {}" v-if="resMsg.msg">
      {{ resMsg.msg }}
    </p>
    <v-icon
      v-if="resMsg.pen"
      name="pr-spinner"
      class="res-msg-spinner anim-spin"
    />
  </div>
</template>

<style lang="scss" scoped>
.res-msg {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--padding-medium);
  gap: var(--padding-medium);
  min-width: 3rem;
  height: fit-content;
  text-align: center;
  filter: drop-shadow(var(--shadow-medium));
  p {
    margin: 0;
    padding: 0;
    font-weight: 600;
  }
  .res-msg-err-icon {
    width: 1.5rem;
    height: 1.5rem;
    color: red;
    fill: red;
  }
  .res-msg-spinner {
    width: 2.5rem;
    height: 2.5rem;
  }
}
</style>
