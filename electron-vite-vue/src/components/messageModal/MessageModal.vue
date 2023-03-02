<script lang="ts" setup>
import { IResMsg } from "../../interfaces/GeneralInterfaces";
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
      <div v-show="msg.pen || msg.err || msg.err" class="icons">
        <v-icon
          v-if="msg.err"
          class="icons-err-icon"
          name="icons-error-round"
        />
        <v-icon
          v-if="msg.pen"
          name="pr-spinner"
          class="icons-spinner anim-spin"
        />
      </div>
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
.icons {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--padding-medium);
  gap: var(--padding-medium);
  min-width: 3rem;
  min-height: 3rem;
  filter: drop-shadow(var(--shadow-medium));
  .icons-err-icon {
    width: 1.5rem;
    height: 1.5rem;
  }
  .icons-spinner {
    width: 2.5rem;
    height: 2.5rem;
  }
}
</style>
