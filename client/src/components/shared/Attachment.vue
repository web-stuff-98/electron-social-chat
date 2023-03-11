<script lang="ts" setup>
import { toRefs, ref, onMounted, onBeforeUnmount } from "vue";
import {
  IResMsg,
  IAttachmentMetadata,
} from "../../interfaces/GeneralInterfaces";
import { attachmentStore } from "../../store/AttachmentStore";
import ResMsg from "../layout/ResMsg.vue";

const props = defineProps<{
  meta?: IAttachmentMetadata;
  msgId: string;
  reverse?: boolean;
}>();
const { meta, msgId } = toRefs(props);

const resMsg = ref<IResMsg>({ msg: "", err: false, pen: false });
const containerRef = ref<HTMLCanvasElement | null>(null);

const observer = new IntersectionObserver(([entry]) => {
  if (entry.isIntersecting) {
    attachmentStore.attachmentEnteredView(msgId.value);
  } else {
    attachmentStore.attachmentLeftView(msgId.value);
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
    ref="containerRef"
    :style="reverse ? { justifyContent: 'flex-end', textAlign: 'right' } : {}"
    class="attachment"
  >
    <div v-if="meta">
      {{ meta }}
    </div>
    <ResMsg :resMsg="resMsg" />
  </div>
</template>

<style lang="scss">
.attachment {
  display: flex;
  flex-direction: row;
  width: 100%;
  justify-content: flex-start;
  text-align: left;
  padding: var(--padding-medium);
}
</style>
