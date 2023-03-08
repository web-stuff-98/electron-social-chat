import { reactive } from "vue";
import {
  IDirectMessage,
  IFriendRequest,
  IInvitation,
} from "../interfaces/GeneralInterfaces";

interface IMessagingStore {
  conversations: {
    uid: string;
    messages: IDirectMessage[];
    friend_requests: IFriendRequest[];
    invitations: IInvitation[];
  }[];
  currentConversationUid: string;
}

export const messagingStore: IMessagingStore = reactive({
  conversations: [],
  currentConversationUid: "",
});
