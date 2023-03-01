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

export const deleteRoom = () =>
  makeRequest("/api/room/delete", {
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

export const getRooms = () =>
  makeRequest("/api/rooms", {
    withCredentials: true,
  });
