import { makeRequest } from "./makeRequest";

export const searchUsers = (username: string) =>
  makeRequest("/api/user/search", {
    withCredentials: true,
    method: "POST",
    data: { username },
  });
