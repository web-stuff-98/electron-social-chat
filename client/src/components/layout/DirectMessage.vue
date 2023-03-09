<script lang="ts" setup>
import { ref, toRefs } from "vue";
import { IDirectMessage } from "../../interfaces/GeneralInterfaces";
import { authStore } from "../../store/AuthStore";
import { messagingStore } from "../../store/MessagingStore";
import { socketStore } from "../../store/SocketStore";
import User from "../shared/User.vue";
const props = defineProps<{ msg: IDirectMessage }>();

const { msg } = toRefs(props);

const isEditing = ref(false);
const editInput = ref("");
const editInputRef = ref<HTMLCanvasElement | null>();

function deleteMessageClicked(id: string) {
  socketStore.send(
    JSON.stringify({
      event_type: "DIRECT_MESSAGE_DELETE",
      ID: id,
      recipient: messagingStore.currentConversationUid,
    })
  );
}

function handleEditClicked() {
  editInput.value = msg.value.content;
  //@ts-ignore
  editInputRef.value = msg.value.content;
  isEditing.value = true;
}

function handleEditInput(e: Event) {
  const target = e.target as HTMLInputElement;
  if (!target || !target.value || target.value.length > 300) return;
  editInput.value = target.value;
}

function handleSubmitUpdate() {
  socketStore.send(
    JSON.stringify({
      event_type: "DIRECT_MESSAGE_UPDATE",
      content: editInput.value,
      recipient: messagingStore.currentConversationUid,
      ID: msg.value.ID,
    })
  );
  isEditing.value = false;
}
</script>

<template>
  <form
    @submit.prevent="handleSubmitUpdate"
    :class="msg.author === authStore.user?.ID ? 'message' : 'message-reversed'"
  >
    <div
      :style="
        msg.author !== authStore.user?.ID ? { justifyContent: 'flex-end' } : {}
      "
      class="user-and-buttons"
    >
      <User
        :menu="msg.author !== authStore.user?.ID"
        :reverse="msg.author !== authStore.user?.ID"
        :small="true"
        :square="true"
        :dateTime="new Date(msg.created_at)"
        :uid="msg.author"
      />
      <div v-if="msg.author === authStore.user?.ID" class="buttons">
        <!-- Delete button -->
        <button
          v-show="!isEditing"
          @click="deleteMessageClicked(msg.ID)"
          type="button"
        >
          <v-icon name="md-delete-sharp" />
        </button>
        <!-- Edit button -->
        <button @click="handleEditClicked" v-show="!isEditing" type="button">
          <v-icon name="ri-edit-box-fill" />
        </button>
      </div>
    </div>
    <p v-show="!isEditing">
      {{ msg.content }}
    </p>
    <div class="edit-textarea-and-buttons">
      <textarea
        autofocus="true"
        ref="editInputRef"
        maxlength="300"
        :value="editInput"
        @input="handleEditInput"
        v-show="isEditing"
      />
      <div class="buttons">
        <!-- Submit update button -->
        <button type="submit" v-show="isEditing">
          <v-icon name="md-send" />
        </button>
        <!-- Cancel update button -->
        <button type="button" @click="isEditing = false" v-show="isEditing">
          <v-icon name="io-close" />
        </button>
      </div>
    </div>
  </form>
</template>

<style lang="scss" scoped>
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
    width: 100%;
    font-size: 0.8rem;
    padding: 0 var(--padding-medium);
    box-sizing: border-box;
    margin: 0;
  }
  .user-and-buttons {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
  }
}
.edit-textarea-and-buttons {
  display: flex;
  align-items: center;
  box-sizing: border-box;
  justify-content: center;
  max-width: 100%;
  gap: 7px;
  padding: var(--padding-medium);
  textarea {
    box-sizing: border-box;
    max-width: calc(var(--aside-width) - var(--padding-medium) * 10 - 2px);
  }
}
.buttons {
  padding-right: 2px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  button {
    display: flex;
    align-items: center;
    height: fit-content;
    padding: 2px;
    svg {
      width: 0.75rem;
      height: 0.75rem;
    }
  }
}
.message-reversed {
  align-items: flex-end;
  text-align: right;
  p {
    text-align: right;
  }
}
</style>
