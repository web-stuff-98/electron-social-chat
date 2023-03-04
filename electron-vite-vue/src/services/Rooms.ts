import { makeRequest } from "./makeRequest";

export const createRoom = (name: string, isPrivate: boolean) =>
  makeRequest("/api/room/create", {
    method: "POST",
    withCredentials: true,
    data: { name, is_private: isPrivate },
  });

export const updateRoom = (id: string, name: string, isPrivate: boolean) =>
  makeRequest(`/api/room/update/${id}`, {
    method: "PATCH",
    withCredentials: true,
    data: { name, is_private: isPrivate },
  });

export const deleteRoom = (id: string) =>
  makeRequest(`/api/room/delete/${id}`, {
    method: "DELETE",
    withCredentials: true,
  });

export const uploadRoomImage = (file: File, id: string) => {
  const data = new FormData();
  data.append("file", file);
  return makeRequest(`/api/room/image/${id}`, {
    method: "POST",
    withCredentials: true,
    data,
  });
};

export const getRooms = (own: boolean) =>
  makeRequest("/api/rooms" + (own ? "?own" : ""), {
    withCredentials: true,
  });

export const getRoom = (id: string) =>
  makeRequest(`/api/room/${id}`, {
    method: "GET",
    withCredentials: true,
  });

export const getRoomDisplayData = (id: string) =>
  makeRequest(`/api/room/display/${id}`, {
    method: "GET",
    withCredentials: true,
  });

export const getRoomChannelsDisplayData = (id: string) =>
  makeRequest(`/api/room/channels/${id}`, {
    method: "GET",
    withCredentials: true,
  });

export const getRoomChannelData = (id: string, roomId: string) =>
  makeRequest(`/api/room/channel/${roomId}/${id}`, {
    method: "GET",
    withCredentials: true,
  });

export const updateRoomChannelsData = (
  roomId: string,
  nameUpdates: { ID: string; name: string }[],
  insertUpdates: { name: string; promote_to_main: boolean }[],
  deleteChannels: string[],
  promote_to_main?: string
) =>
  makeRequest(`/api/room/channels/update/${roomId}`, {
    withCredentials: true,
    method: "PATCH",
    data: {
      update_data: nameUpdates,
      insert_data: insertUpdates,
      delete_ids: deleteChannels,
      promote_to_main: promote_to_main,
    },
  });
