<script lang="ts" setup>
import { getRooms } from "../../../../services/Rooms";
import { roomStore } from "../../../../store/RoomStore";
import { IRoomCard, IResMsg } from "../../../../interfaces/GeneralInterfaces";
import Room from "../../../shared/RoomCard.vue";
import { onMounted, ref } from "vue";
import ResMsg from "../../ResMsg.vue";

const roomsResult = ref<string[]>([]);
const resMsg = ref<IResMsg>({ msg: "", err: false, pen: false });

onMounted(async () => {
  try {
    resMsg.value = { msg: "", err: false, pen: true };
    const rooms: IRoomCard[] = await getRooms();
    roomsResult.value = rooms.map((r) => r.ID);
    roomStore.rooms = rooms;
    resMsg.value = { msg: "", err: false, pen: false };
  } catch (e) {
    resMsg.value = { msg: `${e}`, err: true, pen: true };
  }
});
</script>

<template>
  <div class="container">
    <div class="results-container">
      <div class="room-list">
        <div class="room-container" v-for="id in roomsResult">
          <Room :id="id" />
        </div>
        <div v-if="resMsg.pen || resMsg.err" class="resMsg-container">
          <ResMsg :resMsg="resMsg" />
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.container {
  width: 100%;
  height: 100%;
  .results-container {
    position: relative;
    border: 1px solid var(--base-light);
    border-radius: var(--border-radius-medium);
    overflow: hidden;
    width: 100%;
    height: 100%;
    flex-grow: 1;
    box-sizing: border-box;
    box-shadow: var(--shadow-medium);
    .resMsg-container {
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    .room-list {
      position: absolute;
      width: 100%;
      height: 100%;
      min-width: 100%;
      overflow-y: auto;
      display: flex;
      flex-direction: column;
      justify-content: flex-start;
      align-items: flex-start;
    }
    .room-container {
        width:100%;
    }
  }
}
</style>
