import { reactive } from "vue";
import { IResMsg } from "../interfaces/GeneralInterfaces";

export enum EModalType {
  MESSAGE = "Message",
  LOGIN = "Login",
  REGISTER = "Register",
  WELCOME = "Welcome",
}

// If confirmationCallback or cancellationCallback are undefined then don't show the button
interface IMessageModalProps extends IResMsg {
  confirmationCallback?: Function;
  cancellationCallback?: Function;
}

interface IModalStore {
  showModal: boolean;
  modalType: EModalType;
  messageModalProps: IMessageModalProps;
}

export const modalStore: IModalStore = reactive({
  showModal: true,
  modalType: EModalType.WELCOME,
  messageModalProps: {
    msg: "",
    err: false,
    pen: false,
    confirmationCallback: () => {},
    cancellationCallback: () => {},
  },
});
