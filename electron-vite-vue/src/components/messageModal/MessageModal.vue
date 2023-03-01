<script lang="ts" setup>
import ResMsg from "../layout/ResMsg.vue";
import { ref } from "vue";
import { IResMsg } from "../../interfaces/GeneralInterfaces";
const msg = ref<IResMsg>({ msg: "", err: false, pen: false });
defineProps<{
  confirmationCallback: Function;
  cancellationCallback: Function;
  show: boolean;
}>();
</script>

<template>
  <div v-if="show" class="modal-backdrop" />
  <div v-if="show" class="modal-container">
    <div class="modal">
      <ResMsg :showAnyway="true" :resMsg="msg" />
      <div v-if="!msg.pen" class="buttons">
        <button @click="() => cancellationCallback()">Cancel</button>
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
</style>
