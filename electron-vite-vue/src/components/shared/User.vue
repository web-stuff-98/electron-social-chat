<script lang="ts" setup>
import { toRefs, onMounted, ref, onBeforeUnmount } from "vue";
import { userStore } from "../../store/UserStore";
import useUser from "../../composables/useUser";
const props = defineProps<{ uid: string }>();
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
  <div ref="containerRef" class="user">
    <div
      v-bind:style="{ backgroundImage: `url(${user?.base64pfp})` }"
      class="pfp"
    >
      <v-icon v-if="!user?.base64pfp" name="la-user" />
    </div>
    {{ user?.username }}
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
}
</style>
