<script lang="ts" setup>
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
      <div v-show="msg.pen || msg.err || msg.err" class="msg">
        <v-icon v-if="msg.err" class="msg-err-icon" name="md-error-round" />
        <p v-if="msg.msg">{{ msg.msg }}</p>
        <v-icon
          v-if="msg.pen"
          name="pr-spinner"
          class="msg-spinner anim-spin"
        />
      </div>
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
.msg {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--padding-medium);
  gap: var(--padding-medium);
  min-width: 3rem;
  min-height: 3rem;
  filter: drop-shadow(var(--shadow-medium));
  p {
    margin: 0;
  }
  .msg-err-icon {
    width: 1.5rem;
    height: 1.5rem;
  }
  .msg-spinner {
    width: 2.5rem;
    height: 2.5rem;
  }
}
</style>
