import { createApp } from "vue";
import "./style.css";
import App from "./App.vue";
import "./samples/node-api";

import { OhVueIcon, addIcons } from "oh-vue-icons";
import {
  IoClose,
  FaMinus,
  MdManageaccountsSharp,
  CoMenu,
  PrSpinner,
  MdErrorRound,
  LaUser
} from "oh-vue-icons/icons";

addIcons(
  IoClose,
  FaMinus,
  MdManageaccountsSharp,
  CoMenu,
  PrSpinner,
  MdErrorRound,
  LaUser
);

createApp(App)
  .component("v-icon", OhVueIcon)
  .mount("#app")
  .$nextTick(() => {
    postMessage({ payload: "removeLoading" }, "*");
  });
