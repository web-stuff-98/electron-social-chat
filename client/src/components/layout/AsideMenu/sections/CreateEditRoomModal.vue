<script lang="ts" setup>
import {
  editRoomId,
  showCreateRoom,
  showEditRoom,
} from "../../../../store/CreateEditRoomStore";
import EditRoom from "./EditRoom.vue";
import CreateRoom from "./CreateRoom.vue";
</script>

<template>
  <div v-show="showCreateRoom || showEditRoom" class="modal-backdrop" />
  <div v-show="showCreateRoom || showEditRoom" class="modal-container">
    <div
      v-show="showCreateRoom || showEditRoom"
      class="modal"
      :style="{ padding: '0', paddingTop: 'calc(var(--padding-medium) + 9px)' }"
    >
      <button
        @click="
          {
            showCreateRoom = false;
            showEditRoom = false;
          }
        "
        class="modal-close-button"
      >
        <v-icon name="io-close" />
      </button>
      <EditRoom v-if="editRoomId && showEditRoom" />
      <CreateRoom v-if="showCreateRoom" />
    </div>
  </div>
</template>

<style lang="scss" scoped></style>
