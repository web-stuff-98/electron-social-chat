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
  // Not included from data sent out by server... has v url param which increments whenever the image is updated, so that the image can be refreshed
  img_url?: string;
}
export interface IRoomMessage {
  ID: string;
  content: string;
  author: string;
  created_at: string;
  updated_at: string;
  has_attachment: boolean;
}
export interface IRoomChannel {
  ID: string;
  messages?: IRoomMessage[];
  name: string;
}
export interface IRoom extends IRoomCard {
  channels: string[];
  main_channel: string;
  is_private: boolean;
  members: string[];
  banned: string[];
}
export interface IDirectMessage {
  ID: string;
  content: string;
  author: string;
  created_at: string;
  updated_at: string;
  has_attachment: boolean;
}
export interface IInvitation {
  ID: string;
  author: string;
  recipient: string;
  room_id: string;
  created_at: string;
  accepted: boolean;
  declined: boolean;
}
export interface IFriendRequest {
  ID: string;
  author: string;
  recipient: string;
  created_at: string;
  accepted: boolean;
  declined: boolean;
}
export interface IAttachmentMetadata {
  ID: string;
  meta: string;
  name: string;
  size: number;
  ratio: number;
  failed: boolean;
}
