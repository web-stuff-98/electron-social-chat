<script lang="ts" setup>
import { authStore } from "../../../../store/AuthStore";
import { ref, onMounted, onBeforeUnmount, watch } from "vue";
import User from "../../../shared/User.vue";
import {
  IResMsg,
  IDirectMessage,
  IInvitation,
  IFriendRequest,
} from "../../../../interfaces/GeneralInterfaces";
import { getConversations, getConversation } from "../../../../services/Users";
import ResMsg from "../../ResMsg.vue";

const currentConversationUid = ref("");
const conversations = ref<string[]>([]);
const resMsg = ref<IResMsg>({ msg: "", err: false, pen: false });

type Conversation = {
  messages: IDirectMessage[];
  invitations: IInvitation[];
  friend_requests: IFriendRequest[];
};

const conversation = ref<Conversation>();

watch(currentConversationUid, async (_, newVal) => {
  if (!newVal) return;
  try {
    resMsg.value = { msg: "", err: false, pen: true };
    const data: Conversation = await getConversation(newVal);
    conversation.value = data;
    resMsg.value = { msg: "", err: false, pen: false };
  } catch (e) {
    resMsg.value = { msg: `${e}`, err: false, pen: false };
  }
});

onMounted(async () => {
  const abortController = new AbortController();
  try {
    resMsg.value = { msg: "", err: false, pen: true };
    const uids = await getConversations();
    conversations.value = uids;
    resMsg.value = { msg: "", err: false, pen: false };
  } catch (e) {
    resMsg.value = { msg: `${e}`, err: true, pen: false };
  }
  return () => {
    abortController.abort();
  };
});
</script>

<template>
  <div class="container">
    <div class="messaging-container">
      <div class="messages">
        <div class="message">
          <User
            :small="true"
            :dateTime="new Date()"
            :uid="authStore.user?.ID!"
          />
          <p>
            Lorem ipsum dolor sit amet consectetur adipisicing elit. Iusto vitae
            commodi adipisci corporis, doloribus aperiam odit possimus
            aspernatur saepe eaque.
          </p>
        </div>
      </div>
      <form>
        <input maxlength="300" type="text" />
        <button>
          <v-icon name="md-send" />
        </button>
      </form>
    </div>
    <div class="buttons">
      <button @click="currentConversationUid = ''" type="button">
        All conversations
      </button>
    </div>
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
    padding: var(--padding-medium);
    box-sizing: border-box;
    border: 1px solid var(--base-light);
    border-radius: var(--border-radius-medium);
    form {
      display: flex;
      gap: var(--padding-medium);
      align-items: center;
      box-sizing: border-box;
      padding: 0;
      margin: auto;
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
        svg {
          height: 1.25rem;
          width: 1.25rem;
          filter: drop-shadow(var(--shadow-medium));
        }
      }
    }
    .messages {
      overflow-y: auto;
      flex-grow: 1;
      width: 100%;
      .message,
      .message-reversed {
        width: 100%;
        height: 100%;
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
  .buttons {
    display: flex;
    gap: var(--padding-medium);
    width: 100%;
    button {
      flex-grow: 1;
      padding: var(--padding-medium) var(--padding);
    }
  }
}
</style>
