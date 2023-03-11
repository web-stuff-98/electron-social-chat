<script lang="ts" setup>
import { toRefs, ref, onMounted, onBeforeUnmount } from "vue";
import {
  IResMsg,
  IAttachmentMetadata,
} from "../../interfaces/GeneralInterfaces";
import { attachmentStore } from "../../store/AttachmentStore";
import ResMsg from "../layout/ResMsg.vue";
import ProgressBar from "../shared/Progress.vue";

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
    :style="reverse ? { justifyContent: 'flex-end', textAlign: 'right' } : {}"
    ref="containerRef"
    class="attachment"
  >
    <div v-if="meta && meta.ratio === 1">
      {{ meta }}
    </div>
    <ProgressBar v-if="meta && meta.ratio < 1" :ratio="meta.ratio" />
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
}
</style>
