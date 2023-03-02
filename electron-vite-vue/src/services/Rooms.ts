import { makeRequest } from "./makeRequest";

export const createRoom = (name: string, isPrivate: boolean) =>
  makeRequest("/api/room/create", {
    method: "POST",
    withCredentials: true,
    data: { name, is_private: isPrivate },
  });

export const updateRoom = (name: string, isPrivate: boolean) =>
  makeRequest("/api/room/update", {
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
