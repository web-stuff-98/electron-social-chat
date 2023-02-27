export interface IResMsg {
  msg: string;
  err: boolean;
  pen: boolean;
}
export interface IInputFileEvent extends Event {
  target: HTMLInputElement;
}
