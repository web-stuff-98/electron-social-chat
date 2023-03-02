<script lang="ts" setup>
import { toRefs } from "vue";
import useRoomChannel from "../../composables/useRoomChannel";
const props = defineProps<{ id: string; mainChannelId: string }>();
const { id } = toRefs(props);
const channel = useRoomChannel(id.value);
</script>

<template>
  <div class="channel">
    <div class="name">
      {{ channel?.name }} {{ mainChannelId === id ? "(main)" : "" }}
    </div>
    <div class="buttons">
      <button>
        <v-icon name="md-delete-sharp" />
      </button>
      <button>
        <v-icon name="ri-edit-box-fill" />
      </button>
      <button v-if="mainChannelId !== id">
        <v-icon name="bi-caret-up-fill" />
      </button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.channel {
  display: flex;
  justify-content: space-between;
  gap: var(--padding);
  width: 100%;
  box-sizing: border-box;
  align-items: center;
  padding: 2px;
  .name {
    text-align: left;
    padding: 0 3px;
    flex-grow: 1;
  }
  .buttons {
    border-radius: var(--border-radius-medium);
    gap: 1px;
    display: flex;
    button {
      border: 1px solid var(--base);
      margin: 0;
      padding: 0;
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: 4px;
      svg {
        width: 1rem;
        height: 1rem;
        padding: 1px;
      }
    }
  }
}
</style>
