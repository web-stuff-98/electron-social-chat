import { reactive } from "vue";

interface IPendingCall {
  caller: string;
  called: string;
}

export const pendingCallsStore: IPendingCall[] = reactive([]);
