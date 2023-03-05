<script setup lang="ts">
import {
  EWelcomeModalType,
  welcomeModalStore,
} from "../../store/WelcomeModalStore";
import Welcome from "./sections/Welcome.vue";
import Login from "./sections/Login.vue";
import Register from "./sections/Register.vue";
import { watch } from "vue";

watch(welcomeModalStore, (oldVal, newVal) => {
  if (newVal.modalType !== oldVal.modalType) {
    welcomeModalStore.messageModalProps = {
      ...welcomeModalStore.messageModalProps,
      msg: "",
      err: false,
      pen: false,
    };
  }
});
</script>

<template>
  <div class="modal-backdrop" />
  <div class="modal-container">
    <div class="modal">
      <button
        @click="welcomeModalStore.modalType = EWelcomeModalType.WELCOME"
        v-show="welcomeModalStore.modalType !== EWelcomeModalType.WELCOME"
        class="modal-close-button"
      >
        <v-icon name="io-close" />
      </button>
      <Welcome
        v-if="welcomeModalStore.modalType === EWelcomeModalType.WELCOME"
      />
      <Login v-if="welcomeModalStore.modalType === EWelcomeModalType.LOGIN" />
      <Register
        v-if="welcomeModalStore.modalType === EWelcomeModalType.REGISTER"
      />
    </div>
  </div>
</template>
