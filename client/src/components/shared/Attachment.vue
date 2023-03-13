<script lang="ts" setup>
import { toRefs, ref, onMounted, onBeforeUnmount } from "vue";
import {
  IResMsg,
  IAttachmentMetadata,
} from "../../interfaces/GeneralInterfaces";
import { baseURL } from "../../services/makeRequest";
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
      <img
        :src="`${baseURL}/api/attachment/${msgId}`"
        v-if="
          meta.meta === 'image/jpeg' ||
          meta.meta === 'image/png' ||
          meta.meta === 'image/avif'
        "
      />
      <button :style="{ flexDirection: 'row-reverse' }" type="button" v-else>
        <v-icon name="fa-download" />
        {{
          meta.name.length > 24 ? meta.name.slice(0, 24 - 1) + "..." : meta.name
        }}
      </button>
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
  img {
    border-radius: var(--border-radius-medium);
    filter: drop-shadow(var(--shadow-medium));
    max-width: 80%;
    border: 2px solid var(--base-light);
  }
  button {
    padding: 3px;
    display: flex;
    gap: var(--padding-medium);
    font-size: 0.666rem;
    line-height: 1;
    align-items: center;
    padding: var(--padding-medium);
    box-shadow: none;
    border: none;
  }
  button:hover {
    outline: 1px solid var(--base-light);
  }
}
</style>
