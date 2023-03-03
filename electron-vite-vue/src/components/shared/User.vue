<script lang="ts" setup>
import { toRefs, onMounted, ref, onBeforeUnmount } from "vue";
import { userStore } from "../../store/UserStore";
import useUser from "../../composables/useUser";
const props = defineProps<{ uid: string; reverse?: boolean }>();
const { uid } = toRefs(props);
const user = useUser(uid.value);

const containerRef = ref<HTMLCanvasElement | null>(null);

const observer = new IntersectionObserver(([entry]) => {
  if (entry.isIntersecting) {
    userStore.userEnteredView(uid.value);
  } else {
    userStore.userLeftView(uid.value);
  }
});

onMounted(() => {
  observer.observe(containerRef.value!);
});
onBeforeUnmount(() => {
  observer.disconnect();
});
</script>

<template>
  <div
    :style="reverse ? { flexDirection: 'row-reverse' } : {}"
    ref="containerRef"
    class="user"
  >
    <div
      v-bind:style="{ backgroundImage: `url(${user?.base64pfp})` }"
      class="pfp"
    >
      <v-icon v-if="!user?.base64pfp" name="la-user" />
    </div>
    <div class="name-date-time">
      <div class="name">
        {{ user?.username }}
      </div>
      <div class="date-time">
        <span>01/02/23</span>
        <span>12:34PM</span>
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
  font-size: 0.833rem;
  width: fit-content;
  height: fit-content;
  padding: var(--padding-medium);
  box-sizing: border-box;
  .pfp {
    width: 2rem;
    height: 2rem;
    border: 2px solid var(--base);
    background: var(--foreground);
    background-size: cover;
    border-radius: 50%;
    box-shadow: var(--shadow-medium);
    display: flex;
    align-items: center;
    justify-content: center;
    svg {
      fill: white;
      width: 60%;
      height: 60%;
    }
  }
  .name-date-time {
    padding: 0;
    display: flex;
    flex-direction: column;
    .date-time {
      display: flex;
      flex-direction: column;
      span {
        padding: 0;
        font-size: 0.666rem;
      }
    }
  }
}
</style>
