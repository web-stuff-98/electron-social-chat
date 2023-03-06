import { reactive } from "vue";
import { makeRequest } from "../services/makeRequest";
import { authStore, IUser } from "./AuthStore";
import { socketStore } from "./SocketStore";

/**
 * Interval runs inside App.vue which uses the "disappearedUsers" array to
 * clear the user from the "users" cache after the user has been "disappeared"
 * for longer than 30 seconds
 */

type DisappearedUser = {
  uid: string;
  lastSeen: number;
};

interface IUserStore {
  users: IUser[];

  visibleUsers: string[];
  disappearedUsers: DisappearedUser[];

  userEnteredView: (uid: string) => void;
  userLeftView: (uid: string) => void;

  cacheUserData: (uid: string, force?: boolean) => void;
}

export const userStore: IUserStore = reactive({
  users: [],

  visibleUsers: [],
  disappearedUsers: [],

  userEnteredView: (uid: string) => {
    userStore.disappearedUsers = userStore.disappearedUsers.filter(
      (u) => u.uid !== uid
    );
    userStore.visibleUsers = [...userStore.visibleUsers, uid];
    if (uid !== authStore.user?.ID)
      socketStore.send(JSON.stringify({ event_type: "WATCH_USER", ID: uid }));
    userStore.cacheUserData(uid);
  },
  userLeftView: (uid: string) => {
    const i = userStore.visibleUsers.findIndex((u) => u === uid);
    if (i === -1) return;
    userStore.visibleUsers.splice(i, 1);
    if (userStore.visibleUsers.findIndex((u) => u === uid) === -1) {
      socketStore.send(
        JSON.stringify({ event_type: "STOP_WATCHING_USER", ID: uid })
      );
    }
  },

  cacheUserData: async (uid: string, force?: boolean) => {
    const found = userStore.users.find((u) => u.ID === uid);
    if (found && !force) return;
    try {
      const u = await makeRequest(`/api/user/${uid}`, {
        method: "GET",
      });
      userStore.users = [...userStore.users.filter((u) => u.ID !== uid), u];
    } catch (e) {
      console.warn(`Failed to cache user data for ${uid}: ${e}`);
    }
  },
});
