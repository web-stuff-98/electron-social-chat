<script lang="ts" setup>
import { useRoute } from "vue-router";
import { IRoom, IResMsg } from "../../interfaces/GeneralInterfaces";
import ResMsg from "../layout/ResMsg.vue";
import { ref, onMounted, onBeforeUnmount } from "vue";
import { getRoom } from "../../services/Rooms";

const route = useRoute();
const room = ref<IRoom | undefined>();
const resMsg = ref<IResMsg>({ msg: "", err: false, pen: false });

onMounted(async () => {
  try {
    resMsg.value = { msg: "", err: false, pen: true };
    const data = await getRoom(route.params.id as string);
    room.value = data;
    resMsg.value = { msg: "", err: false, pen: false };
  } catch (e) {
    resMsg.value = { msg: `${e}`, err: true, pen: false };
  }
});

onBeforeUnmount(async () => {});
</script>

<template>
  <div class="container">
    <div
      v-if="!resMsg.pen && !resMsg.err && !resMsg.msg"
      class="channels-container"
    >
      <div class="channels">
        <div class="channel-container" v-for="channelId in room?.channels">
          <button
            :style="channelId === room?.main_channel ? {} : { 'font-weight': 300 }"
            class="channel"
          >
            # {{ channelId }}
          </button>
        </div>
      </div>
    </div>
    <div class="messaging-container">
      <ResMsg :resMsg="resMsg" />
      <div v-if="!resMsg.pen && !resMsg.err && !resMsg.msg" class="content">
        <div class="header">{{ room?.name }}</div>
      </div>
      <form v-if="!resMsg.pen && !resMsg.err && !resMsg.msg">
        <input type="text" />
        <v-icon name="md-send" />
      </form>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.container {
  width: 100%;
  height: 100%;
  display: flex;
  .channels-container {
    width: fit-content;
    box-sizing: border-box;
    padding-left: calc(1px + var(--padding-medium));
    padding-top: calc(1px + var(--padding-medium));
    padding-bottom: calc(var(--padding-medium) * 2 - 1px);
    .channels {
      box-sizing: border-box;
      border: 1px solid var(--base-light);
      border-radius: var(--border-radius-medium);
      box-shadow: var(--shadow-medium);
      height: 100%;
      width: 100%;
      background: var(--foreground);
      display: flex;
      flex-direction: column;
      padding: var(--padding-medium);
      .channel {
        white-space: nowrap;
        font-size: 0.833rem;
        padding: 2px var(--padding-medium);
      }
    }
  }
  .messaging-container {
    width: 100%;
    height: 100%;
    padding: calc(var(--padding-medium) + 1px);
    padding-bottom: var(--padding-medium);
    padding-right: calc(var(--padding-medium) + 2px);
    box-sizing: border-box;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--padding-medium);
    flex-direction: column;
    .content {
      position: relative;
      box-sizing: border-box;
      flex-grow: 1;
      width: 100%;
      border: 1px solid var(--base-light);
      box-shadow: var(--shadow-medium);
      border-radius: var(--border-radius-medium);
      overflow: hidden;
      background: var(--foreground);
      .header {
        position: absolute;
        top: 0;
        width: 100%;
        height: 2.5rem;
        display: flex;
        align-items: center;
        justify-content: center;
        text-align: center;
        border-bottom: 1px solid var(--base-light);
        background: rgba(0, 0, 0, 0.166);
      }
    }
    form {
      width: 100%;
      display: flex;
      gap: var(--padding-medium);
      align-items: center;
      box-sizing: border-box;
      input {
        flex-grow: 1;
        box-sizing: border-box;
        background: var(--foreground);
      }
      svg {
        height: 2rem;
        widows: 2rem;
        filter: drop-shadow(var(--shadow-medium));
      }
    }
  }
}
</style>
