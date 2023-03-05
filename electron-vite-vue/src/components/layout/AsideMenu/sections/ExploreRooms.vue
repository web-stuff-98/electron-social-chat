<script lang="ts" setup>
import { getRooms } from "../../../../services/Rooms";
import { roomStore } from "../../../../store/RoomStore";
import { IRoomCard, IResMsg } from "../../../../interfaces/GeneralInterfaces";
import Room from "../../../shared/RoomCard.vue";
import { onMounted, ref, toRefs } from "vue";
import ResMsg from "../../ResMsg.vue";

const roomsResult = ref<string[]>([]);
const resMsg = ref<IResMsg>({ msg: "", err: false, pen: false });

const props = defineProps<{ own: boolean }>();

const { own } = toRefs(props);

onMounted(async () => {
  try {
    resMsg.value = { msg: "", err: false, pen: true };
    const rooms: IRoomCard[] = await getRooms(own.value);
    roomsResult.value = rooms.map((r) => r.ID);
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
      <div class="pagination-controls">
        <button>
          <v-icon name="bi-caret-left" />
        </button>
        <div class="page-count">
          <span>1/20</span>
          <span>200</span>
        </div>
        <button>
          <v-icon name="bi-caret-right" />
        </button>
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
    .pagination-controls {
      width: 100%;
      height: 2.5rem;
      background: rgba(0, 0, 0, 0.333);
      border-top: 1px solid var(--base-light);
      box-shadow: 0px 0px 3px rgba(0, 0, 0, 0.5);
      display: flex;
      align-items: center;
      justify-content: center;
      position: absolute;
      bottom: 0;
      left:0;
      .page-count {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        text-align: center;
        font-size: 0.8rem;
        padding-top: 2px;
        span {
          padding: 0;
        }
      }
      button {
        padding: 0;
        border: none;
        background: none;
        box-shadow: none;
        display: flex;
        align-items: center;
        svg {
          width: 1rem;
          height: 1rem;
        }
      }
    }
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
      height: calc(100% - 2.5rem);
      min-width: 100%;
      overflow-y: auto;
      display: flex;
      flex-direction: column;
      justify-content: flex-start;
      align-items: flex-start;
    }
    .room-container {
      width: 100%;
    }
  }
}
</style>
