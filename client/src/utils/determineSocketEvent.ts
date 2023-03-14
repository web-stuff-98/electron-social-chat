export type ChangeData = {
  TYPE: "CHANGE";
  METHOD: SocketEventChangeMethodData;
  ENTITY: SocketEventChangeEntityType;
  DATA: { ID: string };
};

export type RoomMessageData = Omit<
  {
    content: string;
    ID: string;
    author: string;
    has_attachment: boolean;
  },
  "TYPE"
>;

export type RoomMessageUpdateData = Omit<
  {
    content: string;
    ID: string;
  },
  "TYPE"
>;

export type ResponseMessageData = Omit<
  {
    DATA: {
      msg: string;
      err: boolean;
    };
  },
  "TYPE"
>;

export type RoomMessageDeleteData = Omit<{ ID: string }, "TYPE">;

export type DirectMessageData = Omit<
  {
    ID: string;
    content: string;
    author: string;
    recipient: string;
    has_attachment: boolean;
  },
  "TYPE"
>;

export type DirectMessageUpdateData = Omit<
  {
    ID: string;
    content: string;
    author: string;
    recipient: string;
  },
  "TYPE"
>;

export type DirectMessageDeleteData = Omit<
  {
    ID: string;
    author: string;
    recipient: string;
  },
  "TYPE"
>;

export type RoomInvitationData = Omit<
  {
    ID: string;
    author: string;
    recipient: string;
    room_id: string;
  },
  "TYPE"
>;

export type RoomInvitationDeleteData = Omit<
  {
    ID: string;
    author: string;
    recipient: string;
  },
  "TYPE"
>;

export type RoomInvitationResponseData = Omit<
  {
    ID: string;
    author: string;
    recipient: string;
    accept: boolean;
  },
  "TYPE"
>;

export type FriendRequestData = Omit<
  {
    ID: string;
    author: string;
    recipient: string;
  },
  "TYPE"
>;

export type FriendRequestDeleteData = Omit<
  {
    ID: string;
    author: string;
    recipient: string;
  },
  "TYPE"
>;

export type FriendRequestResponseData = Omit<
  {
    ID: string;
    author: string;
    recipient: string;
    accept: boolean;
  },
  "TYPE"
>;

export type AttachmentRequestData = Omit<
  {
    ID: string;
    is_room: boolean;
  },
  "TYPE"
>;

export type AttachmentProgressData = Omit<
  {
    ID: string;
    ratio: number;
    err: boolean;
  },
  "TYPE"
>;

export type AttachmentMetadata = Omit<
  {
    ID: string;
    meta: string;
    name: string;
    size: number;
  },
  "TYPE"
>;

export type BannedData = Omit<
  {
    banned: string;
    banner: string;
    room_id: string;
  },
  "TYPE"
>;

export type UnbannedData = BannedData;

export type BlockedData = Omit<
  {
    blocker: string;
  },
  "TYPE"
>;

export type UnblockedData = BlockedData;

export type CallAcknowledgeData = Omit<
  {
    called: string;
    caller: string;
  },
  "TYPE"
>;

export type CallResponseData = Omit<
  {
    called: string;
    caller: string;
    accept: boolean;
  },
  "TYPE"
>;

export type CallLeftData = Omit<{}, "TYPE">;

export type CallWebRTCOfferFromInitiator = Omit<
  {
    signal: string;
    um_stream_id: string;
    dm_stream_id: string;
  },
  "TYPE"
>;

export type CallWebRTCAnswerFromRecipient = Omit<
  {
    signal: string;
    um_stream_id: string;
    dm_stream_id: string;
  },
  "TYPE"
>;

export type CallWebRTCRecipientRequestedReInitialization = Omit<{}, "TYPE">;

export type SocketEventChangeMethodData =
  | "UPDATE"
  | "INSERT"
  | "DELETE"
  | "UPDATE_IMAGE";

export type SocketEventChangeEntityType = "ROOM" | "USER";

export function instanceOfChangeData(object: any): object is ChangeData {
  return object.TYPE === "CHANGE";
}

export function instanceOfRoomMessageData(
  object: any
): object is RoomMessageData {
  return object.TYPE === "OUT_ROOM_MESSAGE";
}
export function instanceOfRoomMessageUpdateData(
  object: any
): object is RoomMessageUpdateData {
  return object.TYPE === "OUT_ROOM_MESSAGE_UPDATE";
}
export function instanceOfRoomMessageDeleteData(
  object: any
): object is RoomMessageUpdateData {
  return object.TYPE === "OUT_ROOM_MESSAGE_DELETE";
}
export function instanceOfDirectMessageData(
  object: any
): object is DirectMessageData {
  return object.TYPE === "OUT_DIRECT_MESSAGE";
}
export function instanceOfDirectMessageUpdateData(
  object: any
): object is DirectMessageUpdateData {
  return object.TYPE === "OUT_DIRECT_MESSAGE_UPDATE";
}
export function instanceOfDirectMessageDeleteData(
  object: any
): object is DirectMessageDeleteData {
  return object.TYPE === "OUT_DIRECT_MESSAGE_DELETE";
}
export function instanceOfRoomInvitationData(
  object: any
): object is RoomInvitationData {
  return object.TYPE === "OUT_ROOM_INVITATION";
}
export function instanceOfRoomInvitationResponseData(
  object: any
): object is RoomInvitationResponseData {
  return object.TYPE === "OUT_ROOM_INVITATION_RESPONSE";
}
export function instanceOfRoomInvitationDeleteData(
  object: any
): object is RoomInvitationDeleteData {
  return object.TYPE === "OUT_ROOM_INVITATION_DELETE";
}
export function instanceOfFriendRequestData(
  object: any
): object is FriendRequestData {
  return object.TYPE === "OUT_FRIEND_REQUEST";
}
export function instanceOfFriendRequestDeleteData(
  object: any
): object is FriendRequestDeleteData {
  return object.TYPE === "OUT_FRIEND_REQUEST_DELETE";
}
export function instanceOfFriendRequestResponseData(
  object: any
): object is FriendRequestResponseData {
  return object.TYPE === "OUT_FRIEND_REQUEST_RESPONSE";
}
export function instanceOfResponseMessageData(
  object: any
): object is ResponseMessageData {
  return object.TYPE === "RESPONSE_MESSAGE";
}
export function instanceOfAttachmentRequestData(
  object: any
): object is AttachmentRequestData {
  return object.TYPE === "ATTACHMENT_REQUEST";
}
export function instanceOfAttachmentProgressData(
  object: any
): object is AttachmentProgressData {
  return object.TYPE === "ATTACHMENT_PROGRESS";
}
export function instanceOfAttachmentMetadata(
  object: any
): object is AttachmentMetadata {
  return object.TYPE === "ATTACHMENT_METADATA";
}
export function instanceOfBanData(object: any): object is BannedData {
  return object.TYPE === "BANNED";
}
export function instanceOfUnBanData(object: any): object is BannedData {
  return object.TYPE === "UNBANNED";
}
export function instanceOfBlockData(object: any): object is BannedData {
  return object.TYPE === "BLOCKED";
}
export function instanceOfUnblockData(object: any): object is BannedData {
  return object.TYPE === "UNBLOCKED";
}
export function instanceOfCallAcknowledgeData(
  object: any
): object is CallAcknowledgeData {
  return object.TYPE === "CALL_USER_ACKNOWLEDGE";
}
export function instanceOfCallResponseData(
  object: any
): object is CallResponseData {
  return object.TYPE === "CALL_USER_RESPONSE";
}
export function instanceOfCallLeftData(object: any): object is CallLeftData {
  return object.TYPE === "CALL_LEFT";
}
export function instanceOfCallWebRTCOfferFromInitiator(
  object: any
): object is CallWebRTCOfferFromInitiator {
  return object.TYPE === "CALL_WEBRTC_OFFER_FROM_INITIATOR";
}
export function instanceOfCallWebRTCAnswerFromRecipient(
  object: any
): object is CallWebRTCAnswerFromRecipient {
  return object.TYPE === "CALL_WEBRTC_ANSWER_FROM_RECIPIENT";
}
export function instanceOfCallWebRTCRecipientRequestedReInitialization(
  object: any
): object is CallWebRTCRecipientRequestedReInitialization {
  return object.TYPE === "CALL_WEBRTC_REQUESTED_REINITIALIZATION";
}

export function parseSocketEventData(e: MessageEvent): object | undefined {
  let data = JSON.parse(e.data);
  if (!data["DATA"]) return data;
  data["DATA"] = JSON.parse(data["DATA"]);
  return data;
}
