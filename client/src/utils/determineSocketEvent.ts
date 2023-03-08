export type ChangeData = {
  TYPE: "CHANGE";
  METHOD: SocketEventChangeMethod;
  ENTITY: SocketEventChangeEntityType;
  DATA: { ID: string };
};

export type RoomMessageData = Omit<
  {
    content: string;
    ID: string;
    author: string;
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

export type SocketEventChangeMethod =
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
export function instanceOfResponseMessageData(
  object: any
): object is ResponseMessageData {
  return object.TYPE === "RESPONSE_MESSAGE";
}

export function parseSocketEventData(e: MessageEvent): object | undefined {
  let data = JSON.parse(e.data);
  if (!data["DATA"]) return data;
  data["DATA"] = JSON.parse(data["DATA"]);
  return data;
}
