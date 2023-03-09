<script lang="ts" setup>
import { authStore } from "../../../../store/AuthStore";
import {
  ref,
  onMounted,
  onBeforeUnmount,
  watch,
  watchEffect,
  nextTick,
} from "vue";
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
import DirectMessage from "../../DirectMessage.vue";
import {
  instanceOfDirectMessageData,
  parseSocketEventData,
} from "../../../../utils/determineSocketEvent";
import FriendRequest from "../../FriendRequest.vue";
import Invitation from "../../Invitation.vue";

const resMsg = ref<IResMsg>({ msg: "", err: false, pen: false });

enum EListSection {
  "MESSAGING" = "Messaging",
  "FRIEND_REQUESTS" = "Friend requests",
  "INVITATIONS" = "Invitations",
}
const listSection = ref(EListSection.MESSAGING);
const listBottomRef = ref<HTMLCanvasElement | null>();

onMounted(async () => {
  const abortController = new AbortController();
  messagingStore.currentConversationUid = "";
  try {
    resMsg.value = { msg: "", err: false, pen: true };
    const {
      conversations: uids,
      friend_requests,
      invitations,
    }: {
      conversations: string[];
      friend_requests: IFriendRequest[];
      invitations: IInvitation[];
    } = await getConversations();
    uids.forEach((uid) => {
      userStore.cacheUserData(uid, true);
    });
    messagingStore.conversations = uids.map((uid) => ({
      uid,
      messages: [],
    }));
    messagingStore.friend_requests = friend_requests;
    messagingStore.invitations = invitations;
    resMsg.value = { msg: "", err: false, pen: false };
  } catch (e) {
    resMsg.value = { msg: `${e}`, err: true, pen: false };
  }

  socketStore.socket?.addEventListener("message", messageEventListener);

  return () => {
    abortController.abort();
  };
});

watch(listBottomRef, async () => {
  await nextTick(() => {
    // @ts-ignore
    listBottomRef.value?.scrollIntoView({ behavior: "auto" });
  });
});

async function messageEventListener(e: MessageEvent) {
  const data = parseSocketEventData(e);
  if (!data) return;
  if (instanceOfDirectMessageData(data)) {
    if (
      data.author === authStore.user?.ID ||
      data.author === messagingStore.currentConversationUid
    ) {
      if (listBottomRef) {
        await nextTick(() => {
          listBottomRef.value?.scrollIntoView({ behavior: "auto" });
        });
      }
    }
  }
}

onBeforeUnmount(() => {
  socketStore.socket?.removeEventListener("message", messageEventListener);
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
      recipient: messagingStore.currentConversationUid,
    })
  );
  // @ts-ignore
  msgInputRef.value = "";
  msgInput.value = "";
}

async function openConv(uid: string) {
  const abortController = new AbortController();
  try {
    resMsg.value = { msg: "", err: false, pen: true };
    const messages: IDirectMessage[] = await getConversation(uid);
    const convI = messagingStore.conversations.findIndex((c) => c.uid === uid);
    if (convI === -1)
      messagingStore.conversations.push({
        uid,
        messages,
      });
    else if (uid !== "")
      messagingStore.conversations[convI] = {
        uid,
        messages,
      };
    messagingStore.currentConversationUid = uid;
    resMsg.value = { msg: "", err: false, pen: false };
  } catch (e) {
    resMsg.value = { msg: `${e}`, err: false, pen: false };
  }
  return () => {
    abortController.abort();
  };
}
</script>

<template>
  <div class="container">
    <div class="messaging-container">
      <div class="friend-requests-invitations-buttons">
        <button
          @click="listSection = EListSection.MESSAGING"
          v-if="listSection !== EListSection.MESSAGING"
          :style="{ width: '100%', fontSize: '0.75rem' }"
        >
          Back to direct messages
        </button>
        <button
          @click="listSection = EListSection.FRIEND_REQUESTS"
          v-if="listSection === EListSection.MESSAGING"
          :style="{ width: '50%' }"
        >
          Friend requests
        </button>
        <button
          @click="listSection = EListSection.INVITATIONS"
          v-if="listSection === EListSection.MESSAGING"
          :style="{ width: '50%' }"
        >
          Room invitations
        </button>
      </div>
      <div
        v-if="
          !messagingStore.currentConversationUid &&
          listSection === EListSection.MESSAGING
        "
        class="users"
      >
        <button
          @click="() => openConv(uid)"
          class="user"
          v-for="{ uid } in messagingStore.conversations"
        >
          <User :uid="uid" />
        </button>
      </div>
      <div v-if="messagingStore.currentConversationUid" class="list-container">
        <div v-if="listSection === EListSection.MESSAGING" class="list">
          <DirectMessage
            v-for="msg in messagingStore.conversations.find(
              (c) => c.uid === messagingStore.currentConversationUid
            )?.messages || []"
            :msg="msg"
          />
        </div>
        <div v-if="listSection === EListSection.FRIEND_REQUESTS" class="list">
          <FriendRequest
            :frq="frq"
            v-for="frq in messagingStore.friend_requests"
          />
        </div>
        <div v-if="listSection === EListSection.INVITATIONS" class="list">
          <Invitation :inv="inv" v-for="inv in messagingStore.invitations" />
        </div>
        <div ref="listBottomRef" class="list-bottom" />
      </div>
      <form
        class="message-form"
        @submit.prevent="handleFormSubmit"
        v-if="
          messagingStore.currentConversationUid &&
          listSection === EListSection.MESSAGING
        "
      >
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
      class="all-conversations-button"
      v-if="
        messagingStore.currentConversationUid &&
        listSection === EListSection.MESSAGING
      "
      @click="messagingStore.currentConversationUid = ''"
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
  .friend-requests-invitations-buttons {
    display: flex;
    width: 100%;
    button:nth-child(2) {
      border-left: 1px solid var(--base-light);
    }
    button {
      font-size: 0.7rem;
      border: none;
      border-bottom: 2px solid var(--base-light);
      border-radius: 0;
    }
  }
  .messaging-container {
    display: flex;
    flex-direction: column;
    width: 100%;
    flex-grow: 1;
    box-sizing: border-box;
    border: 2px solid var(--base);
    border-radius: var(--border-radius-medium);
    box-shadow: var(--shadow-medium);
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
        border: none;
        border-radius: 0;
        background: none;
        box-shadow: none;
      }
      .user:hover {
        background: var(--foreground-hover);
      }
    }
    .message-form {
      display: flex;
      gap: var(--padding-medium);
      align-items: center;
      box-sizing: border-box;
      width: 100%;
      padding: 3px;
      box-sizing: border-box;
      border-top: 1px solid var(--base-light);
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
    .list-container {
      flex-grow: 1;
      width: 100%;
      position: relative;
      box-sizing: border-box;
      .list {
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
      }
    }
  }
  button {
    width: 100%;
    padding: var(--padding-medium) var(--padding);
  }
  .all-conversations-button {
    border: 2px solid var(--base-light);
  }
}
</style>
