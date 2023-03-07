<script lang="ts" setup>
import { ref } from "vue";
import { IResMsg } from "../../../../interfaces/GeneralInterfaces";
import { searchUsers } from "../../../../services/Users";
import { userStore } from "../../../../store/UserStore";

import ResMsg from "../../ResMsg.vue";
import User from "../../../shared/User.vue";

/*
Debounced search using setTimeout
*/

const usernameSearchInput = ref("");
const usersResult = ref<string[]>([]);
const searchTimeout = ref<NodeJS.Timeout>();

const resMsg = ref<IResMsg>({ msg: "", err: false, pen: false });

async function search() {
  const abortController = new AbortController();
  try {
    const users: string[] = await searchUsers(usernameSearchInput.value);
    usersResult.value = users;
    users.forEach((uid) => userStore.cacheUserData(uid));
    resMsg.value = { msg: "", err: false, pen: false };
  } catch (e) {
    resMsg.value = { msg: `${e}`, err: true, pen: false };
  }
  return () => {
    abortController.abort();
  };
}

function handleUsernameInput(e: Event) {
  const target = e.target as HTMLInputElement;
  if (target.value.length > 16) {
    resMsg.value = { msg: "Max 16 characters", err: true, pen: false };
    return;
  }
  resMsg.value = { msg: "", err: false, pen: true };
  usernameSearchInput.value = target.value;
  if (searchTimeout.value) clearTimeout(searchTimeout.value);
  searchTimeout.value = setTimeout(search, 300);
}
</script>

<template>
  <div class="container">
    <div class="results-container">
      <div class="user-list">
        <div class="user-container" v-for="user in usersResult">
          <User :menu="true" :uid="user" />
        </div>
      </div>
    </div>
    <form>
      <input id="username" @input="handleUsernameInput" type="text" />
      <v-icon class="no-fill" v-if="!resMsg.pen" name="hi-search" />
      <v-icon v-if="resMsg.pen" class="anim-spin" name="pr-spinner" />
    </form>
    <ResMsg v-if="!resMsg.pen" :resMsg="resMsg" />
  </div>
</template>

<style lang="scss" scoped>
.container {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  align-items: center;
  box-sizing: border-box;
  gap: var(--padding-medium);
  .results-container {
    position: relative;
    border: 1px solid var(--base-light);
    border-radius: var(--border-radius-medium);
    width: 100%;
    flex-grow: 1;
    box-sizing: border-box;
    .user-list {
      position: absolute;
      width: 100%;
      height: 100%;
      overflow-y: auto;
      display: flex;
      flex-direction: column;
      justify-content: flex-start;
      align-items: flex-start;
    }
  }
  form {
    box-sizing: border-box;
    display: flex;
    gap: 2px;
    align-items: center;
    max-width: 100%;
    width: 100%;
    padding: 0;
    margin: 0;
    input {
      margin: 0;
      box-sizing: border-box;
      width: calc(100% - 1.5rem);
    }
    svg {
      width: 1.5rem;
      height: 1.5rem;
    }
    .no-fill {
      fill: none;
    }
  }
}
</style>
