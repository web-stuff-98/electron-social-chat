<script lang="ts" setup>
import { getRoomPage } from "../../../../services/Rooms";
import { IRoomCard, IResMsg } from "../../../../interfaces/GeneralInterfaces";
import Room from "../../RoomCard.vue";
import { onMounted, onBeforeUnmount, ref, toRefs } from "vue";
import ResMsg from "../../ResMsg.vue";
import { socketStore } from "../../../../store/SocketStore";
import {
  instanceOfChangeData,
  instanceOfRoomMessageDeleteData,
  parseSocketEventData,
} from "../../../../utils/determineSocketEvent";

const roomsResult = ref<string[]>([]);
const page = ref<number>(1);
const resMsg = ref<IResMsg>({ msg: "", err: false, pen: false });
const props = defineProps<{ own: boolean }>();

const { own } = toRefs(props);

async function getPage() {
  try {
    resMsg.value = { msg: "", err: false, pen: true };
    const rooms: IRoomCard[] = await getRoomPage(page.value, own.value);
    if (!rooms || rooms.length === 0) {
      roomsResult.value = rooms.map((r) => r.ID);
      resMsg.value = { msg: "No results", err: false, pen: false };
    } else {
      roomsResult.value = rooms.map((r) => r.ID);
      resMsg.value = { msg: "", err: false, pen: false };
    }
  } catch (e) {
    resMsg.value = { msg: `${e}`, err: true, pen: true };
  }
}

function watchForDeletes(e: MessageEvent) {
  const data = parseSocketEventData(e);
  if (!data) return;
  if (instanceOfChangeData(data)) {
    if (data.ENTITY === "ROOM" && data.METHOD === "DELETE") {
      const i = roomsResult.value.findIndex((r) => r === data.DATA.ID);
      if (i !== -1) {
        roomsResult.value.splice(i, 1);
      }
    }
  }
}
onMounted(() => {
  getPage();
  socketStore.socket?.addEventListener("message", watchForDeletes);
});
onBeforeUnmount(() => {
  socketStore.socket?.removeEventListener("message", watchForDeletes);
});
</script>

<template>
  <div class="container">
    <div class="results-container">
      <div class="room-list">
        <div class="room-container" v-for="id in roomsResult">
          <Room :id="id" />
        </div>
        <div
          v-if="resMsg.pen || resMsg.err || resMsg.msg"
          class="resMsg-container"
        >
          <ResMsg :resMsg="resMsg" />
        </div>
      </div>
      <div class="pagination-controls">
        <button
          @click="
            {
              if (page !== 1) {
                page--;
                getPage();
              }
            }
          "
        >
          <v-icon name="bi-caret-left" />
        </button>
        <div class="page-count">{{ page }}</div>
        <button
          @click="
            {
              if (roomsResult && roomsResult.length > 0) {
                page++;
                getPage();
              }
            }
          "
        >
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
    border: 2px solid var(--base);
    border-radius: var(--border-radius-medium);
    overflow: hidden;
    box-shadow: var(--shadow-medium);
    width: 100%;
    height: 100%;
    flex-grow: 1;
    box-sizing: border-box;
    .pagination-controls {
      width: 100%;
      height: 1.5rem;
      background: rgba(0, 0, 0, 0.333);
      border-top: 1px solid var(--base-light);
      display: flex;
      align-items: center;
      justify-content: center;
      position: absolute;
      bottom: 0;
      left: 0;
      .page-count {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        text-align: center;
        font-size: 0.8rem;
        padding-top: 2px;
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
      height: calc(100% - 1.5rem);
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
