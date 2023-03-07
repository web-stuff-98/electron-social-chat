<script lang="ts" setup>
import { authStore } from "../../../../store/AuthStore";
import User from "../../../shared/User.vue";
</script>

<template>
  <div class="container">
    <div class="message">
      <User :small="true" :dateTime="new Date()" :uid="authStore.user?.ID!" />
      <p>
        Lorem ipsum dolor sit amet consectetur adipisicing elit. Iusto vitae
        commodi adipisci corporis, doloribus aperiam odit possimus aspernatur
        saepe eaque.
      </p>
    </div>
  </div>
</template>
<style lang="scss" scoped>
.container {
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  height: 100%;
  width: 100%;
  padding: var(--padding-medium);
  box-sizing: border-box;
  top: 0;
  left: 0;
}
.message,
.message-reversed {
  width: 100%;
  height: 100%;
  margin: 0;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  box-sizing: border-box;
  text-align: left;
  p {
    text-align: left;
    font-size: 0.6rem;
    padding: 0 var(--padding-medium);
    margin: 0;
  }
}
.message-reversed {
  align-items: flex-end;
  text-align: right;
  p {
    text-align: right;
  }
}
</style>
