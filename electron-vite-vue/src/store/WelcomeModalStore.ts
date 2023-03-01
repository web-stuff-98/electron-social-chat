import { reactive } from "vue";
import { IResMsg } from "../interfaces/GeneralInterfaces";

export enum EWelcomeModalType {
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
  modalType: EWelcomeModalType;
  messageModalProps: IMessageModalProps;
}

export const welcomeModalStore: IModalStore = reactive({
  showModal: true,
  modalType: EWelcomeModalType.WELCOME,
  messageModalProps: {
    msg: "",
    err: false,
    pen: false,
    confirmationCallback: () => {},
    cancellationCallback: () => {},
  },
});
