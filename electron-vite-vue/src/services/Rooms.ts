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