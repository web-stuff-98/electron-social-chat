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

export function parseSocketEventData(e: MessageEvent): object | undefined {
  let data = JSON.parse(e.data);
  if (!data["DATA"]) return data;
  data["DATA"] = JSON.parse(data["DATA"]);
  return data;
}
