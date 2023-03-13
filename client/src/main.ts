import { createApp, watch } from "vue";
import { createRouter, createWebHashHistory } from "vue-router";
import "./style.css";
import App from "./App.vue";
import Home from "./components/routes/Home.vue";
import Call from "./components/routes/Call.vue";
import Room from "./components/routes/Room.vue";
import EditRoom from "./components/layout/AsideMenu/sections/EditRoom.vue";

import { OhVueIcon, addIcons } from "oh-vue-icons";
import {
  IoClose,
  FaMinus,
  BiShieldFill,
  CoMenu,
  PrSpinner,
  MdErrorRound,
  FaUser,
  HiSearch,
  MdDeleteSharp,
  BiDoorClosedFill,
  RiEditBoxFill,
  MdSend,
  BiCaretUpFill,
  BiCaretDownFill,
  IoAddCircleSharp,
  BiCheckLg,
  BiCaretLeft,
  BiCaretRight,
  MdDarkmode,
  MdAttachfileRound,
  FaDownload,
  HiPhoneMissedCall,
  BiMicFill,
  BiMicMuteFill,
  MdScreenshare,
  BiCameraVideo,
  MdStopscreenshare,
  BiCameraVideoOff,
  HiPhoneIncoming,
} from "oh-vue-icons/icons";
import { authStore } from "./store/AuthStore";

addIcons(
  IoClose,
  FaMinus,
  BiShieldFill,
  CoMenu,
  PrSpinner,
  MdErrorRound,
  FaUser,
  HiSearch,
  MdDeleteSharp,
  BiDoorClosedFill,
  RiEditBoxFill,
  MdSend,
  BiCaretUpFill,
  BiCaretDownFill,
  IoAddCircleSharp,
  BiCheckLg,
  BiCaretLeft,
  BiCaretRight,
  MdDarkmode,
  MdAttachfileRound,
  FaDownload,
  HiPhoneMissedCall,
  HiPhoneIncoming,
  BiMicFill,
  BiMicMuteFill,
  MdScreenshare,
  MdStopscreenshare,
  BiCameraVideo,
  BiCameraVideoOff
);

const routes = [
  { path: "/", component: Call, name: "home" },
  { path: "/room/:id", component: Room, name: "room" },
  { path: "/room/edit/:id", component: EditRoom, name: "edit_room" },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

router.beforeEach(async (to) => {
  if (!authStore.user && to.path !== "/") {
    router.push("/");
    return "/";
  }
});

watch(authStore, (_, newVal) => {
  if (!newVal.user) {
    router.push("/");
  }
});

createApp(App)
  .use(router)
  .component("v-icon", OhVueIcon)
  .mount("#app")
  .$nextTick(() => {
    postMessage({ payload: "removeLoading" }, "*");
  });
