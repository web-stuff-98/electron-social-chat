export type ChangeData = {
  TYPE: "CHANGE";
  METHOD: SocketEventChangeMethod;
  ENTITY: SocketEventChangeEntityType;
  DATA: { ID: string };
};

export type SocketEventChangeMethod =
  | "UPDATE"
  | "INSERT"
  | "DELETE"
  | "UPDATE_IMAGE";

export type SocketEventChangeEntityType = "ROOM" | "USER";

export function instanceOfChangeData(object: any): object is ChangeData {
  return object.TYPE === "CHANGE";
}

export function parseSocketEventData(e: MessageEvent): object | undefined {
  let data = JSON.parse(e.data);
  if (!data["DATA"]) return;
  data["DATA"] = JSON.parse(data["DATA"]);
  return data
}
