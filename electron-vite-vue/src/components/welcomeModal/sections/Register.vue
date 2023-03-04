<script setup lang="ts">
import { ref } from "vue";
import { authStore } from "../../../store/AuthStore";
import ResMsg from "../../layout/ResMsg.vue";

const username = ref("");
const password = ref("");
</script>

<template>
  <form @submit.prevent="authStore.register(username, password)">
    <label for="username">Username</label>
    <input autofocus v-model="username" id="username" type="text" />
    <label for="password">Password</label>
    <input v-model="password" id="password" type="password" />
    <button type="submit" v-if="!authStore.resMsg.pen">Create account</button>
    <ResMsg :resMsg="authStore.resMsg" />
  </form>
</template>

<style lang="scss" scoped>
form {
  display: flex;
  flex-direction: column;
  button {
    margin-top: var(--padding);
  }
  input {
    text-align: center;
  }
  label,
  input {
    margin-bottom: var(--padding-medium);
  }
}
</style>
