import { reactive } from "vue";
import { IResMsg } from "../interfaces/GeneralInterfaces";
import { makeRequest } from "../services/makeRequest";
import { EModalType, modalStore } from "./ModalStore";
import { socketStore } from "./SocketStore";

export interface IUser {
  ID: string;
  username: string;
  base64pfp: string;
  online: boolean;
}

export interface IAuthStore {
  user?: IUser;
  login: (username: string, password: string) => void;
  register: (username: string, password: string) => void;
  logout: () => void;
  deleteAccount: () => void;
  uploadPfp: (file: File) => Promise<void>;
  refreshToken: () => void;
  resMsg: IResMsg;
}

export const authStore: IAuthStore = reactive({
  user: undefined,
  resMsg: { msg: "", err: false, pen: false },
  login: async (username: string, password: string) => {
    try {
      authStore.resMsg = { msg: "", err: false, pen: true };
      const user = await makeRequest("/api/acc/login", {
        withCredentials: true,
        method: "POST",
        data: { username, password },
      });
      authStore.user = user;
      authStore.resMsg = { msg: "", err: false, pen: false };
      modalStore.showModal = false;
      socketStore.connectSocket(user.ID);
    } catch (e) {
      authStore.resMsg = { msg: `${e}`, err: true, pen: false };
    }
  },
  register: async (username: string, password: string) => {
    try {
      authStore.resMsg = { msg: "", err: false, pen: true };
      const user = await makeRequest("/api/acc/register", {
        withCredentials: true,
        method: "POST",
        data: { username, password },
      });
      authStore.user = user;
      authStore.resMsg = { msg: "", err: false, pen: false };
      modalStore.showModal = false;
      socketStore.connectSocket(user.ID);
    } catch (e) {
      authStore.resMsg = { msg: `${e}`, err: true, pen: false };
    }
  },
  logout: async () => {
    try {
      authStore.resMsg = { msg: "", err: false, pen: true };
      await makeRequest("/api/acc/logout", {
        withCredentials: true,
        method: "POST",
      });
      authStore.user = undefined;
      authStore.resMsg = { msg: "", err: false, pen: false };
      modalStore.modalType = EModalType.WELCOME;
      modalStore.showModal = true;
    } catch (e) {
      authStore.resMsg = { msg: `${e}`, err: true, pen: false };
    }
  },
  deleteAccount: async () => {
    try {
      authStore.resMsg = { msg: "", err: false, pen: true };
      await makeRequest("/api/acc/delete", {
        withCredentials: true,
        method: "DELETE",
      });
      authStore.user = undefined;
      authStore.resMsg = { msg: "", err: false, pen: false };
      modalStore.modalType = EModalType.WELCOME;
      modalStore.showModal = true;
    } catch (e) {
      authStore.resMsg = { msg: `${e}`, err: true, pen: false };
    }
  },
  uploadPfp: (file: File) => {
    const data = new FormData();
    data.append("file", file);
    return makeRequest("/api/acc/pfp", {
      method: "POST",
      withCredentials: true,
      data,
    });
  },
  refreshToken: async () => {
    try {
      await makeRequest("/api/acc/refresh", {
        withCredentials: true,
        method: "POST",
      });
    } catch (e) {
      authStore.user = undefined;
    }
  },
});
