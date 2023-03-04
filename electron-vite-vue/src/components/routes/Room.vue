<script lang="ts" setup>
import { useRoute } from "vue-router";
import {
  IRoom,
  IResMsg,
  IRoomChannel,
} from "../../interfaces/GeneralInterfaces";
import ResMsg from "../layout/ResMsg.vue";
import { roomChannelStore } from "../../store/RoomChannelStore";
import { roomStore } from "../../store/RoomStore";
import { ref, watch, onMounted, onBeforeUnmount } from "vue";
import { getRoom, getRoomChannelData } from "../../services/Rooms";
import { socketStore } from "../../store/SocketStore";
import RoomMessage from "../layout/RoomMessage.vue";
import { authStore } from "../../store/AuthStore";
import {
  parseSocketEventData,
  instanceOfRoomMessageData,
  instanceOfRoomMessageUpdateData,
  instanceOfRoomMessageDeleteData,
} from "../../utils/determineSocketEvent";

const route = useRoute();
const room = ref<IRoom | undefined>();
const resMsg = ref<IResMsg>({ msg: "", err: false, pen: false });

const messageInput = ref("");
const currentChannel = ref("");

function messageEventListener(e: MessageEvent) {
  const data = parseSocketEventData(e);
  console.log(data);
  if (!data) return;
  const i = roomChannelStore.channels.findIndex(
    (c) => c.ID === room.value?.main_channel
  );
  if (instanceOfRoomMessageData(data)) {
    roomChannelStore.channels[i].messages?.push({
      ID: data.ID,
      author: data.author,
      content: data.content,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    });
    return;
  }
  if (instanceOfRoomMessageUpdateData(data)) {
    const msgI: number =
      roomChannelStore.channels[i].messages?.findIndex(
        (m) => m.ID === data.ID
      ) || -1;
    roomChannelStore.channels[i].messages![msgI] = {
      ...roomChannelStore.channels[i].messages![msgI],
      content: data.content,
      updated_at: new Date().toISOString(),
    };
  }
  if (instanceOfRoomMessageDeleteData(data)) {
    const msgI = roomChannelStore.channels[i].messages?.findIndex(
      (m) => m.ID === data.ID
    );
    roomChannelStore.channels[i].messages?.splice(msgI as number, 1);
  }
}

onMounted(async () => {
  try {
    resMsg.value = { msg: "", err: false, pen: true };
    const data: IRoom = await getRoom(route.params.id as string);
    room.value = data;
    currentChannel.value = data.main_channel;
    roomStore.rooms = [
      ...roomStore.rooms.filter((r) => r.ID !== data.ID),
      data,
    ];
    await roomChannelStore.getDisplayDataForChannels(route.params.id as string);
    await roomChannelStore.getFullDataForChannel(
      data.main_channel,
      room.value.ID
    );
    socketStore.send(
      JSON.stringify({
        event_type: "ROOM_OPEN_CHANNEL",
        channel: currentChannel.value,
      })
    );
    roomStore.currentRoom = data.ID;
    resMsg.value = { msg: "", err: false, pen: false };
  } catch (e) {
    resMsg.value = { msg: `${e}`, err: true, pen: false };
  }
  socketStore.socket?.addEventListener("message", messageEventListener);
});

onBeforeUnmount(async () => {
  roomStore.currentRoom = "";
  socketStore.send(
    JSON.stringify({
      event_type: "ROOM_EXIT_CHANNEL",
      channel: currentChannel.value,
    })
  );
  socketStore.socket?.removeEventListener("message", messageEventListener);
});

watch(currentChannel, async (oldChannel, channel) => {
  if (!channel) return;
  try {
    if (oldChannel) {
      socketStore.send(
        JSON.stringify({
          event_type: "ROOM_EXIT_CHANNEL",
          channel: oldChannel,
        })
      );
    }
    await roomChannelStore.getFullDataForChannel(
      channel,
      room.value?.ID as string
    );
    socketStore.send(
      JSON.stringify({
        event_type: "ROOM_OPEN_CHANNEL",
        channel: channel,
      })
    );
  } catch (e) {
    resMsg.value = {
      msg: `Error retreiving channel data: ${e}`,
      err: true,
      pen: false,
    };
  }
});

function handleMessageInput(e: Event) {
  const target = e.target as HTMLInputElement;
  if (target.value.length > 300 || !target.value.trim()) return;
  messageInput.value = target.value;
}

function handleMessageSubmit() {
  socketStore.send(
    JSON.stringify({
      event_type: "ROOM_MESSAGE",
      content: messageInput.value,
      channel: currentChannel.value,
    })
  );
  messageInput.value = "";
}
</script>

<template>
  <div class="container">
    <div
      v-if="!resMsg.pen && !resMsg.err && !resMsg.msg"
      class="channels-container"
    >
      <div class="channels">
        <!-- Main channel -->
        <div class="channel-container">
          <button
            type="button"
            @click="currentChannel = room?.main_channel!"
            :style="{ fontWeight: 600, marginBottom: 'var(--padding)', ...(currentChannel !== room?.main_channel! ? {
              filter:'opacity(0.5)'
            } : {}) }"
            class="channel"
          >
            #
            {{
              roomChannelStore.channels.find((c) => c.ID === room?.main_channel)
                ?.name
            }}
          </button>
        </div>
        <!-- Secondary channels -->
        <div
          class="channel-container"
          v-for="channel in roomChannelStore.channels.filter(
            (c) => c.ID !== room?.main_channel
          )"
        >
          <button
            @click="currentChannel = channel.ID"
            type="button"
            class="channel"
            :style="
              currentChannel !== channel.ID ? { filter: 'opacity(0.5)' } : {}
            "
          >
            # {{ channel.name }}
          </button>
        </div>
      </div>
    </div>
    <div class="messaging-container">
      <ResMsg :resMsg="resMsg" />
      <div v-if="!resMsg.pen && !resMsg.err && !resMsg.msg" class="content">
        <div class="header">{{ room?.name }}</div>
        <div class="messages-list-container">
          <div class="messages">
            <div
              v-for="message in roomChannelStore.channels.find(
                (c) => c.ID === currentChannel
              )?.messages"
              class="message-container"
            >
              <RoomMessage
                :reverse="message.author !== authStore.user?.ID"
                :msg="message"
              />
            </div>
          </div>
        </div>
      </div>
      <form
        @submit.prevent="handleMessageSubmit"
        v-if="!resMsg.pen && !resMsg.err && !resMsg.msg"
      >
        <input
          maxlength="300"
          :value="messageInput"
          @input="handleMessageInput"
          type="text"
        />
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
      gap: var(--padding-medium);
      .channel {
        white-space: nowrap;
        font-size: 0.833rem;
        padding: 2px var(--padding-medium);
        width: 100%;
        text-align: left;
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
      display: flex;
      flex-direction: column;
      .messages-list-container {
        flex-grow: 1;
        width: 100%;
        height: 100%;
        position: relative;
        .messages {
          display: flex;
          flex-direction: column;
          width: 100%;
          height: 100%;
          overflow-y: auto;
          position: absolute;
          left: 0;
          top: 0;
          .message-container {
            padding: 0;
          }
        }
      }
      .header {
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
