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
  }[];
  invitations: IInvitation[];
  friend_requests: IFriendRequest[];
  currentConversationUid: string;
}

export const messagingStore: IMessagingStore = reactive({
  conversations: [],
  invitations: [],
  friend_requests: [],
  currentConversationUid: "",
});
