<script lang="ts" setup>
import { toRefs, onMounted, ref, onBeforeUnmount } from "vue";
import { userStore } from "../../store/UserStore";
import useUser from "../../composables/useUser";
import { userdropdownStore } from "../../store/UserDropdownStore";
const props = defineProps<{
  uid: string;
  dateTime?: Date;
  reverse?: boolean;
  menu?: boolean;
  room?: string;
  small?: boolean;
  square?: boolean;
}>();
const { uid, menu, room } = toRefs(props);
const user = useUser(uid.value);

const containerRef = ref<HTMLCanvasElement | null>(null);

const observer = new IntersectionObserver(([entry]) => {
  if (entry.isIntersecting) {
    userStore.userEnteredView(uid.value);
  } else {
    userStore.userLeftView(uid.value);
  }
});

function handleClick() {
  if (!menu.value) return;
  userdropdownStore.roomId = room ? room?.value || "" : "";
  userdropdownStore.subject = uid.value;
  userdropdownStore.open = true;
}

onMounted(() => {
  observer.observe(containerRef.value!);
});
onBeforeUnmount(() => {
  observer.disconnect();
});
</script>

<template>
  <div
    @click="handleClick"
    :style="{
      ...(reverse ? { flexDirection: 'row-reverse' } : {}),
      ...(menu ? { cursor: 'pointer' } : {}),
    }"
    ref="containerRef"
    class="user"
  >
    <div
      :style="{
        ...(small ? { width: '1.666rem', height: '1.666rem' } : {}),
        ...(user?.base64pfp
          ? { backgroundImage: `url(${user?.base64pfp})` }
          : {}),
        ...(square ? { borderRadius: 'var(--border-radius-medium)' } : {}),
      }"
      class="pfp"
    >
      <v-icon v-if="!user?.base64pfp" name="fa-user" />
    </div>
    <div class="name-date-time">
      <div
        :style="small ? { fontSize: '0.75rem', fontWeight: 300 } : {}"
        class="name"
      >
        {{ user?.username }}
      </div>
      <div v-if="dateTime" class="date-time">
        <span :style="small ? { fontSize: '0.666rem' } : {}">{{
          new Intl.DateTimeFormat("en-GB").format(dateTime)
        }}</span>
        <span :style="small ? { fontSize: '0.666rem' } : {}">{{
          new Intl.DateTimeFormat("default", {
            hour: "numeric",
            minute: "numeric",
            second: "numeric",
          }).format(dateTime)
        }}</span>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.user {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  font-size: 0.866rem;
  width: fit-content;
  height: fit-content;
  padding: var(--padding-medium);
  box-sizing: border-box;
  .pfp {
    width: 2rem;
    height: 2rem;
    border: 2px solid var(--text-color);
    background: var(--foreground);
    background-size: cover;
    border-radius: 50%;
    box-shadow: var(--shadow-medium);
    display: flex;
    align-items: center;
    justify-content: center;
    svg {
      fill: var(--text-color);
      width: 60%;
      height: 60%;
    }
  }
  .name-date-time {
    padding: 0;
    display: flex;
    flex-direction: column;
    font-weight: 600;
    .date-time {
      display: flex;
      flex-direction: column;
      span {
        font-weight: 300;
        padding: 0;
        font-size: 0.65rem;
      }
    }
  }
}
</style>
