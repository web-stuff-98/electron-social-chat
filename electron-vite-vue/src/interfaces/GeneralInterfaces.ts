export interface IResMsg {
  msg: string;
  err: boolean;
  pen: boolean;
}
export interface IRoomCard {
  ID: string;
  name: string;
  blur: string;
  author: string;
}
export interface IRoomMessage {
  ID: string;
  content: string;
  author: string;
  created_at: string;
  updated_at: string;
}
export interface IRoomChannel {
  ID: string;
  messages?: IRoomMessage[];
  name: string;
}
export interface IRoom extends IRoomCard {
  channels: IRoomChannel;
  main_channel: string;
}
