<script lang="ts" setup>
import { IResMsg } from "../../interfaces/GeneralInterfaces";
import ResMsg from "../layout/ResMsg.vue";
defineProps<{
  confirmationCallback: Function;
  cancellationCallback?: Function;
  show: boolean;
  msg: IResMsg;
}>();
</script>

<template>
  <div v-if="show" class="modal-backdrop" />
  <div v-if="show" class="modal-container">
    <div class="modal">
      <ResMsg :resMsg="msg" />
      <p v-if="msg.msg">{{ msg.msg }}</p>
      <div v-if="!msg.pen" class="buttons">
        <button
          v-if="cancellationCallback"
          @click="() => cancellationCallback!()"
        >
          Cancel
        </button>
        <button @click="() => confirmationCallback()">Confirm</button>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.buttons {
  width: 100%;
  gap: var(--padding-medium);
  display: flex;
  margin-top: var(--padding-medium);
  button {
    flex-grow: 1;
  }
}
p {
  margin: 0;
  text-align: center;
}
</style>
