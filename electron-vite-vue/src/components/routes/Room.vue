<script lang="ts" setup>
import { useRoute } from "vue-router";
import { IRoom, IResMsg } from "../../interfaces/GeneralInterfaces";
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
    <div class="content">
      <div class="header">{{ room?.name }}</div>
    </div>
    <form>
      <input type="text" />
      <v-icon name="md-send" />
    </form>
  </div>
</template>

<style lang="scss" scoped>
.container {
  width: 100%;
  height: 100%;
  padding: calc(var(--padding-medium) + 1px);
  box-sizing: border-box;
  display: flex;
  gap: var(--padding-medium);
  flex-direction: column;
  .content {
    position: relative;
    flex-grow: 1;
    border: 1px solid var(--base-light);
    box-shadow: var(--shadow-medium);
    border-radius: var(--border-radius-medium);
    overflow: hidden;
    .header {
      position: absolute;
      top: 0;
      width: 100%;
      height: 1.75rem;
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
    }
    svg {
      height: 2rem;
      widows: 2rem;
      filter: drop-shadow(var(--shadow-medium));
    }
  }
}
</style>
