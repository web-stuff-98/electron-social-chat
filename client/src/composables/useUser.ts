import { userStore } from "../store/UserStore";
import { authStore, IUser } from "../store/AuthStore";

const useUser = (uid: string): IUser | undefined =>
  authStore.user && uid === authStore.user?.ID
    ? authStore.user
    : userStore.users.find((u) => u.ID === uid);

export default useUser;
