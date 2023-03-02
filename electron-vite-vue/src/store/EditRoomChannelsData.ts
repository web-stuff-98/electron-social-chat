import { reactive } from "vue";

interface IEditRoomChannelsData {
  flaggedForDeletion: string[];
  updateData: { ID: string; name: string }[];
  insertData: { name: string; promoteToMain: boolean }[];
  promoteToMain: string;
}

export const editRoomChannelsData: IEditRoomChannelsData = reactive({
  flaggedForDeletion: [],
  updateData: [],
  insertData: [],
  promoteToMain: "",
});
