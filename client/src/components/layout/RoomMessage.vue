<script lang="ts" setup>
import { IRoomMessage } from "../../interfaces/GeneralInterfaces";
import User from "../shared/User.vue";
import { ref, toRefs } from "vue";
import { roomChannelStore } from "../../store/RoomChannelStore";
import { authStore } from "../../store/AuthStore";
import { socketStore } from "../../store/SocketStore";
import { roomStore } from "../../store/RoomStore";
const props = defineProps<{
  msg: IRoomMessage;
  reverse?: boolean;
}>();

const { msg } = toRefs(props);

const isEditing = ref(false);
const editInput = ref("");
const editInputRef = ref<HTMLCanvasElement | null>(null);

function handleEditInput(e: Event) {
  const target = e.target as HTMLInputElement;
  if (!target || !target.value || target.value.length > 300) return;
  editInput.value = target.value;
}

function handleSubmitUpdate() {
  socketStore.send(
    JSON.stringify({
      event_type: "ROOM_MESSAGE_UPDATE",
      content: editInput.value,
      channel: roomChannelStore.currentChannel,
      id: msg.value.ID,
    })
  );
  isEditing.value = false;
}

function handleDeleteMessage() {
  socketStore.send(
    JSON.stringify({
      event_type: "ROOM_MESSAGE_DELETE",
      channel: roomChannelStore.currentChannel,
      id: msg.value.ID,
    })
  );
}
</script>

<template>
  <form
    @submit.prevent="handleSubmitUpdate"
    :style="reverse ? { textAlign: 'right', flexDirection: 'row-reverse' } : {}"
    class="message"
  >
    <div class="user-container">
      <User
        :small="true"
        :square="true"
        :room="roomStore.currentRoom"
        :menu="msg.author !== authStore.user?.ID"
        :dateTime="new Date(msg.created_at)"
        :reverse="reverse"
        :uid="msg.author"
      />
    </div>
    <div v-show="!isEditing" class="content">
      {{ msg.content }}
    </div>
    <textarea
      autofocus="true"
      ref="editInputRef"
      maxlength="300"
      @value="editInput"
      @input="handleEditInput"
      v-show="isEditing"
    />
    <div class="icons">
      <!-- Submit update button -->
      <button type="submit" v-show="isEditing">
        <v-icon name="md-send" />
      </button>
      <!-- Cancel update button -->
      <button type="button" @click="isEditing = false" v-show="isEditing">
        <v-icon name="io-close" />
      </button>
      <!-- Delete button -->
      <button
        @click="handleDeleteMessage"
        type="button"
        v-show="!isEditing && msg.author === authStore.user?.ID"
        class="delete-button"
      >
        <v-icon name="md-delete-sharp" />
      </button>
      <!-- Edit button -->
      <button
        type="button"
        v-show="!isEditing && msg.author === authStore.user?.ID"
        @click="
          {
            //@ts-ignore
            editInputRef.value = msg.content;
            isEditing = true;
            editInput = msg.content;
          }
        "
        class="edit-button"
      >
        <v-icon name="ri-edit-box-fill" />
      </button>
    </div>
  </form>
</template>

<style lang="scss" scoped>
.message {
  text-align: left;
  font-size: 0.8rem;
  display: flex;
  align-items: center;
  .content {
    padding: 0;
    flex-grow: 1;
    margin: var(--padding-medium);
  }
  textarea {
    flex-grow: 1;
    margin: var(--padding-medium);
  }
  .icons {
    display: flex;
    flex-direction: column;
    gap: var(--padding-medium);
    padding: var(--padding-medium);
    button {
      padding: 0;
      align-items: center;
      justify-content: center;
      display: flex;
      padding: var(--padding-small);
      border: 2px solid var(--base);
      width: 1.5rem;
      height: 1.5rem;
      svg {
        width: 90%;
        height: 90%;
      }
    }
  }
}
</style>
