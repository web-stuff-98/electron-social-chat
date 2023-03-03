import { reactive } from "vue";

interface IEditRoomChannelsData {
  delete_ids: string[];
  update_data: { ID: string; name: string }[];
  insert_data: { name: string; promote_to_main: boolean }[];
  promote_to_main: string;
}

export const editRoomChannelsData: IEditRoomChannelsData = reactive({
  delete_ids: [],
  update_data: [],
  insert_data: [],
  promote_to_main: "",
});
