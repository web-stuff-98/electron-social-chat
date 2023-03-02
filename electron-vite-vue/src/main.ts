import { createApp } from "vue";
import { createRouter, createWebHashHistory } from "vue-router";
import "./style.css";
import App from "./App.vue";
import Home from "./components/routes/Home.vue";
import Room from "./components/routes/Room.vue";
import EditRoom from "./components/routes/EditRoom.vue";
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
  MdSend,
  BiCaretUpFill,
  BiCaretDownFill,
  IoAddCircleSharp,
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
  RiEditBoxFill,
  MdSend,
  BiCaretUpFill,
  BiCaretDownFill,
  IoAddCircleSharp
);

const routes = [
  { path: "/", component: Home, name: "home" },
  { path: "/room/:id", component: Room, name: "room" },
  { path: "/room/edit/:id", component: EditRoom, name: "edit_room" },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

createApp(App)
  .use(router)
  .component("v-icon", OhVueIcon)
  .mount("#app")
  .$nextTick(() => {
    postMessage({ payload: "removeLoading" }, "*");
  });
