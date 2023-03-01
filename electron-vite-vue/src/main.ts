import { createApp } from "vue";
import "./style.css";
import App from "./App.vue";
import "./samples/node-api";

import { OhVueIcon, addIcons } from "oh-vue-icons";
import {
  IoClose,
  FaMinus,
  BiShieldFill,
  CoMenu,
  PrSpinner,
  MdErrorRound,
  LaUser,
  HiSearch,
  MdDeleteSharp,
  BiDoorClosedFill,
  RiEditBoxFill,
} from "oh-vue-icons/icons";

addIcons(
  IoClose,
  FaMinus,
  BiShieldFill,
  CoMenu,
  PrSpinner,
  MdErrorRound,
  LaUser,
  HiSearch,
  MdDeleteSharp,
  BiDoorClosedFill,
  RiEditBoxFill
);

createApp(App)
  .component("v-icon", OhVueIcon)
  .mount("#app")
  .$nextTick(() => {
    postMessage({ payload: "removeLoading" }, "*");
  });
