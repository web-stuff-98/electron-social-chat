import { reactive } from "vue";
import { IAttachmentMetadata } from "../interfaces/GeneralInterfaces";
import { getAttachmentMetadata } from "../services/Attachment";

/**
 * Interval runs inside App.vue which uses the "disappearedAttachments" array to
 * clear the attachment metadata cache after the attacmhent has been "disappeared"
 * for longer than 30 seconds
 */

interface IAttachmentStore {
  attachmentMetadata: IAttachmentMetadata[];

  visibleAttachments: string[];
  disappearedAttachments: DisappearedAttachment[];

  attachmentEnteredView: (id: string) => void;
  attachmentLeftView: (id: string) => void;

  getMetadata: (id: string) => IAttachmentMetadata | undefined;

  cacheAttachmentData: (id: string, force?: boolean) => void;
}

type DisappearedAttachment = {
  id: string;
  lastSeen: number;
};

export const attachmentStore: IAttachmentStore = reactive({
  attachmentMetadata: [],

  visibleAttachments: [],
  disappearedAttachments: [],

  getMetadata: (id: string) =>
    attachmentStore.attachmentMetadata.find((a) => a.ID === id),

  attachmentEnteredView: (id: string) => {
    attachmentStore.disappearedAttachments =
      attachmentStore.disappearedAttachments.filter((a) => a.id !== id);
    attachmentStore.visibleAttachments = [
      ...attachmentStore.visibleAttachments,
      id,
    ];
    attachmentStore.cacheAttachmentData(id);
  },

  attachmentLeftView: (id: string) => {
    const i = attachmentStore.visibleAttachments.findIndex((a) => a === id);
    if (i === -1) return;
    attachmentStore.visibleAttachments.splice(i, 1);
  },

  cacheAttachmentData: async (id: string, force?: boolean) => {
    const found = attachmentStore.attachmentMetadata.find((a) => a.ID === id);
    if (found && !force) return;
    try {
      const a = await getAttachmentMetadata(id);
      attachmentStore.attachmentMetadata.push(a);
    } catch (e) {
      console.warn(`Failed to cache attachment metadata for ${id}: ${e}`);
    }
  },
});
