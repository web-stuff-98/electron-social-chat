<script lang="ts" setup>
import { useRoute } from "vue-router";
import { IResMsg } from "../../interfaces/GeneralInterfaces";
import ResMsg from "../layout/ResMsg.vue";
import { roomChannelStore } from "../../store/RoomChannelStore";
import { roomStore } from "../../store/RoomStore";
import { ref, onMounted, onBeforeUnmount } from "vue";
import { socketStore } from "../../store/SocketStore";
import RoomMessage from "../layout/RoomMessage.vue";
import { authStore } from "../../store/AuthStore";
import {
  parseSocketEventData,
  instanceOfRoomMessageData,
  instanceOfRoomMessageUpdateData,
  instanceOfRoomMessageDeleteData,
  instanceOfAttachmentRequestData,
} from "../../utils/determineSocketEvent";
import MessageModal from "../messageModal/MessageModal.vue";
import { uploadAttachment } from "../../services/Attachment";

const route = useRoute();
const resMsg = ref<IResMsg>({ msg: "", err: false, pen: false });
const messagesBottomRef = ref<HTMLElement>();

const modalConfirmation = ref(() => {});
const modalCancellation = ref<Function | undefined>(() => {});
const showModal = ref(false);
const modalMsg = ref<IResMsg>({ msg: "", err: false, pen: false });

const attachmentFile = ref<File | null>();
const attachmentInputRef = ref<HTMLElement>();
const messageInput = ref("");

function messageEventListener(e: MessageEvent) {
  const data = parseSocketEventData(e);
  if (!data) return;
  const i = roomChannelStore.channels.findIndex(
    (c) => c.ID === roomChannelStore.currentChannel
  );
  if (instanceOfRoomMessageData(data)) {
    roomChannelStore.channels[i].messages?.push({
      ID: data.ID,
      author: data.author,
      content: data.content,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      has_attachment: data.has_attachment,
    });
    if (messagesBottomRef) {
      messagesBottomRef.value?.scrollIntoView({ behavior: "auto" });
    }
    return;
  }
  if (instanceOfRoomMessageUpdateData(data)) {
    const msgI = roomChannelStore.channels[i].messages?.findIndex(
      (m) => m.ID === data.ID
    );
    if (msgI === undefined || msgI === -1) return;
    console.log(roomChannelStore.channels[i].messages![msgI]);
    roomChannelStore.channels[i].messages![msgI].content = data.content;
    roomChannelStore.channels[i].messages![msgI].updated_at =
      new Date().toISOString();
  }
  if (instanceOfRoomMessageDeleteData(data)) {
    const msgI = roomChannelStore.channels[i].messages?.findIndex(
      (m) => m.ID === data.ID
    );
    roomChannelStore.channels[i].messages?.splice(msgI as number, 1);
  }
}

onMounted(async () => {
  const abortController = new AbortController();
  try {
    resMsg.value = { msg: "", err: false, pen: true };
    const data = await roomStore.roomEnteredView(
      route.params.id as string,
      true
    );
    roomChannelStore.currentChannel = data.main_channel as string;
    await roomChannelStore.getDisplayDataForChannels(route.params.id as string);
    await roomChannelStore.getFullDataForChannel(
      data.main_channel as string,
      route.params.id as string
    );
    socketStore.send(
      JSON.stringify({
        event_type: "ROOM_OPEN_CHANNEL",
        channel: roomChannelStore.currentChannel,
      })
    );
    roomStore.currentRoom = route.params.id as string;
    resMsg.value = { msg: "", err: false, pen: false };
  } catch (e) {
    resMsg.value = { msg: `${e}`, err: true, pen: false };
  }

  socketStore.socket?.addEventListener("message", messageEventListener);
  socketStore.socket?.addEventListener("message", watchForAttachmentRequest);

  return () => {
    abortController.abort();
  };
});

onBeforeUnmount(async () => {
  roomStore.currentRoom = "";
  roomStore.roomLeftView(route.params.id as string);
  socketStore.send(
    JSON.stringify({
      event_type: "ROOM_EXIT_CHANNEL",
      channel: roomChannelStore.currentChannel,
    })
  );
  socketStore.socket?.removeEventListener("message", messageEventListener);
  socketStore.socket?.removeEventListener("message", watchForAttachmentRequest);
});

function handleMessageInput(e: Event) {
  const target = e.target as HTMLInputElement;
  if (target.value.length > 300 || !target.value.trim()) return;
  messageInput.value = target.value;
}

function handleFormSubmit() {
  socketStore.send(
    JSON.stringify({
      event_type: "ROOM_MESSAGE",
      content: messageInput.value,
      channel: roomChannelStore.currentChannel,
      has_attachment: attachmentFile.value ? true : false,
    })
  );
  messageInput.value = "";
}

async function uploadRequestedAttachment(msgId: string) {
  const file = attachmentFile.value;
  attachmentFile.value = undefined;
  if (file) {
    try {
      await uploadAttachment(
        file,
        msgId,
        roomChannelStore.currentChannel,
        true
      );
    } catch (e) {
      showModal.value = true;
      modalMsg.value = {
        msg: `${e}`,
        err: true,
        pen: false,
      };
      modalCancellation.value = undefined;
      modalConfirmation.value = () => (showModal.value = false);
    }
  } else {
    showModal.value = true;
    modalMsg.value = {
      msg: `No file selected`,
      err: true,
      pen: false,
    };
    modalCancellation.value = undefined;
    modalConfirmation.value = () => (showModal.value = false);
  }
}

async function openChannel(id: string) {
  if (roomChannelStore.currentChannel === id) return;
  try {
    if (roomChannelStore.currentChannel) {
      socketStore.send(
        JSON.stringify({
          event_type: "ROOM_EXIT_CHANNEL",
          channel: roomChannelStore.currentChannel,
        })
      );
    }
    await roomChannelStore.getFullDataForChannel(id, route.params.id as string);
    roomChannelStore.currentChannel = id;
    socketStore.send(
      JSON.stringify({
        event_type: "ROOM_OPEN_CHANNEL",
        channel: id,
      })
    );
  } catch (e) {
    resMsg.value = {
      msg: `Error retreiving channel data: ${e}`,
      err: true,
      pen: false,
    };
  }
}

function watchForAttachmentRequest(e: MessageEvent) {
  const data = parseSocketEventData(e);
  if (!data) return;
  if (instanceOfAttachmentRequestData(data)) {
    if (data.is_room) {
      uploadRequestedAttachment(data.ID);
    }
  }
}

function selectAttachment(e: Event) {
  const target = e.target as HTMLInputElement;
  if (!target.files || !target.files[0]) return;
  if (target.files[0].size > 20 * 1024 * 1024) {
    modalMsg.value = {
      msg: "File too large. Max 20mb.",
      err: true,
      pen: false,
    };
    showModal.value = true;
    modalCancellation.value = undefined;
    modalConfirmation.value = () => (showModal.value = false);
    return;
  }
  attachmentFile.value = target.files[0];
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
            @click="
              roomChannelStore.currentChannel = roomStore.getRoom(
                route.params.id as string
              )?.main_channel!
            "
            :style="{ fontWeight: 600, marginBottom: 'var(--padding)', ...(roomChannelStore.currentChannel !== roomStore.getRoom(route.params.id as string)?.main_channel! ? {
              filter:'opacity(0.5)'
            } : {}) }"
            class="channel"
          >
            #
            {{
              roomChannelStore.channels.find(
                (c) =>
                  c.ID ===
                  roomStore.getRoom(route.params.id as string)?.main_channel
              )?.name
            }}
          </button>
        </div>
        <!-- Secondary channels -->
        <div
          class="channel-container"
          v-for="channel in roomChannelStore.channels.filter(
            (c) => c.ID !== roomStore.getRoom(route.params.id as string)?.main_channel
          )"
        >
          <button
            @click="openChannel(channel.ID)"
            type="button"
            class="channel"
            :style="
              roomChannelStore.currentChannel !== channel.ID
                ? { filter: 'opacity(0.5)' }
                : {}
            "
          >
            # {{ channel.name }}
          </button>
        </div>
      </div>
    </div>
    <div class="messaging-container">
      <MessageModal
        :msg="modalMsg"
        :show="showModal"
        :confirmationCallback="modalConfirmation"
        :cancellationCallback="modalCancellation"
      />
      <ResMsg :resMsg="resMsg" />
      <div v-if="!resMsg.pen && !resMsg.err && !resMsg.msg" class="content">
        <div class="header">
          {{ roomStore.getRoom(route.params.id as string)?.name }}
        </div>
        <div class="messages-list-container">
          <div class="messages">
            <div
              v-for="message in roomChannelStore.channels.find(
                (c) => c.ID === roomChannelStore.currentChannel
              )?.messages"
              class="message-container"
            >
              <RoomMessage
                :reverse="message.author !== authStore.user?.ID"
                :msg="message"
              />
              <div ref="messagesBottomRef" class="messages-bottom" />
            </div>
          </div>
        </div>
      </div>
      <form
        @submit.prevent="handleFormSubmit"
        v-if="!resMsg.pen && !resMsg.err && !resMsg.msg"
      >
        <input
          maxlength="300"
          :value="messageInput"
          @input="handleMessageInput"
          type="text"
        />
        <button type="submit">
          <v-icon name="md-send" />
        </button>
        <button @click="attachmentInputRef?.click()" type="button">
          <v-icon name="md-attachfile-round" />
        </button>
        <input
          @change="selectAttachment"
          ref="attachmentInputRef"
          type="file"
        />
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
      border: 2px solid var(--base-light);
      border-radius: var(--border-radius-medium);
      box-shadow: var(--shadow-medium);
      height: 100%;
      width: 100%;
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
        border: 2px solid var(--base-light);
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
      border: 2px solid var(--base-light);
      border-radius: var(--border-radius-medium);
      box-shadow: var(--shadow-medium);
      overflow: hidden;
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
        background: rgba(0, 0, 0, 0.1);
        font-weight: 600;
      }
    }
    form {
      width: 100%;
      display: flex;
      gap: var(--padding-medium);
      align-items: center;
      box-sizing: border-box;
      button {
        border: none;
        background: none;
        box-shadow: none;
        padding: 0;
      }
      input {
        flex-grow: 1;
        box-sizing: border-box;
        border: 2px solid var(--base-light);
      }
      svg {
        height: 2rem;
        filter: drop-shadow(var(--shadow-medium));
      }
    }
  }
}
</style>
