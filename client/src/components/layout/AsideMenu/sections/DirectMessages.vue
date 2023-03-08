<script lang="ts" setup>
import { authStore } from "../../../../store/AuthStore";
import { ref, onMounted, onBeforeUnmount, watch } from "vue";
import {
  IResMsg,
  IDirectMessage,
  IInvitation,
  IFriendRequest,
} from "../../../../interfaces/GeneralInterfaces";
import { getConversations, getConversation } from "../../../../services/Users";
import User from "../../../shared/User.vue";
import ResMsg from "../../ResMsg.vue";
import { userStore } from "../../../../store/UserStore";
import { socketStore } from "../../../../store/SocketStore";
import { messagingStore } from "../../../../store/MessagingStore";

const currentConversationUid = ref("");
const resMsg = ref<IResMsg>({ msg: "", err: false, pen: false });

type Conversation = {
  messages: IDirectMessage[];
  invitations: IInvitation[];
  friend_requests: IFriendRequest[];
};

watch(currentConversationUid, async (_, newVal) => {
  if (!newVal) return;
  const abortController = new AbortController();
  try {
    resMsg.value = { msg: "", err: false, pen: true };
    const data: Conversation = await getConversation(newVal);
    const convI = messagingStore.conversations.findIndex(
      (c) => c.uid === newVal
    );
    if (convI === -1)
      messagingStore.conversations.push({
        uid: newVal,
        ...data,
      });
    else if (newVal !== "")
      messagingStore.conversations[convI] = {
        uid: newVal,
        ...data,
      };
    resMsg.value = { msg: "", err: false, pen: false };
  } catch (e) {
    resMsg.value = { msg: `${e}`, err: false, pen: false };
  }
  return () => {
    abortController.abort();
  };
});

onMounted(async () => {
  const abortController = new AbortController();
  try {
    resMsg.value = { msg: "", err: false, pen: true };
    const uids: string[] = await getConversations();
    uids.forEach((uid) => {
      userStore.cacheUserData(uid);
    });
    messagingStore.conversations = uids.map((uid) => ({
      uid,
      messages: [],
      invitations: [],
      friend_requests: [],
    }));
    resMsg.value = { msg: "", err: false, pen: false };
  } catch (e) {
    resMsg.value = { msg: `${e}`, err: true, pen: false };
  }
  return () => {
    abortController.abort();
  };
});

const msgInput = ref("");
const msgInputRef = ref<HTMLCanvasElement | null>();
function handleFormInput(e: Event) {
  const target = e.target as HTMLInputElement;
  if (!target || !target.value || target.value.length > 300) return;
  msgInput.value = target.value;
  // @ts-ignore
  msgInputRef.value = target.value;
}
function handleFormSubmit() {
  if (!msgInput.value || msgInput.value.length > 300) return;
  socketStore.send(
    JSON.stringify({
      event_type: "DIRECT_MESSAGE",
      content: msgInput.value,
      recipient: currentConversationUid.value,
    })
  );
  // @ts-ignore
  msgInputRef.value = "";
  msgInput.value = "";
}
</script>

<template>
  <div class="container">
    <div class="messaging-container">
      <div v-if="!currentConversationUid" class="users">
        <button
          @click="currentConversationUid = uid"
          class="user"
          v-for="{ uid } in messagingStore.conversations"
        >
          <User :uid="uid" />
        </button>
      </div>
      <div v-if="currentConversationUid" class="messages-container">
        <div class="messages-list">
          <div
            v-for="msg in messagingStore.conversations.find(
              (c) => c.uid === currentConversationUid
            )?.messages || []"
            :class="
              msg.author === authStore.user?.ID ? 'message' : 'message-reversed'
            "
          >
            <User
              :reverse="msg.author !== authStore.user?.ID"
              :small="true"
              :dateTime="new Date(msg.created_at)"
              :uid="msg.author"
            />
            <p>
              {{ msg.content }}
            </p>
          </div>
        </div>
      </div>
      <form @submit.prevent="handleFormSubmit" v-if="currentConversationUid">
        <input
          ref="msgInputRef"
          @input="handleFormInput"
          :value="msgInput"
          maxlength="300"
          type="text"
        />
        <button type="submit">
          <v-icon name="md-send" />
        </button>
      </form>
    </div>
    <button
      v-if="currentConversationUid"
      @click="currentConversationUid = ''"
      type="button"
    >
      All conversations
    </button>
    <ResMsg :resMsg="resMsg" />
  </div>
</template>
<style lang="scss" scoped>
.container {
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  height: 100%;
  width: 100%;
  padding: var(--padding);
  gap: var(--padding-medium);
  box-sizing: border-box;
  top: 0;
  left: 0;
  .messaging-container {
    display: flex;
    flex-direction: column;
    width: 100%;
    flex-grow: 1;
    padding: 3px;
    box-sizing: border-box;
    border: 1px solid var(--base-light);
    border-radius: var(--border-radius-medium);
    overflow: hidden;
    .users {
      display: flex;
      flex-direction: column;
      gap: var(--padding-medium);
      width: 100%;
      box-sizing: border-box;
      overflow-y: auto;
      .user {
        padding: 0;
        cursor: pointer;
        width: 100%;
        padding: 3px;
        display: flex;
        justify-content: flex-start;
      }
    }
    form {
      display: flex;
      gap: var(--padding-medium);
      align-items: center;
      box-sizing: border-box;
      padding: 0;
      width: 100%;
      input {
        flex-grow: 1;
        box-sizing: border-box;
        background: var(--foreground);
        width: calc(100% - var(--padding-medium));
      }
      button {
        display: flex;
        align-items: center;
        border: none;
        background: none;
        box-shadow: none;
        padding: 0;
        margin: 0;
        width: fit-content;
        svg {
          height: 1.25rem;
          width: 1.25rem;
          filter: drop-shadow(var(--shadow-medium));
        }
      }
    }
    .messages-container {
      flex-grow: 1;
      width: 100%;
      position: relative;
      .messages-list {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
        position: absolute;
        left: 0;
        top: 0;
        overflow-y: auto;
        box-sizing: border-box;
        padding-bottom: var(--padding-medium);
        .message,
        .message-reversed {
          width: 100%;
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
      }
    }
  }
  button {
    width: 100%;
    padding: var(--padding-medium) var(--padding);
  }
}
</style>
