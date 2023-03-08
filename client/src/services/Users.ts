import { makeRequest } from "./makeRequest";

export const searchUsers = (username: string) =>
  makeRequest("/api/user/search", {
    withCredentials: true,
    method: "POST",
    data: { username },
  });

export const getConversation = (uid: string) =>
  makeRequest(`/api/acc/conversation/${uid}`, {
    withCredentials: true,
  });

export const getConversations = () =>
  makeRequest("/api/acc/conversations", {
    withCredentials: true,
  });
