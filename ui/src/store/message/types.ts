export enum MessageActionTypes {
  SHOW_MESSAGE = '@@message/SHOW'
}

export interface MessageState {
  readonly text: string | null
  readonly isError: boolean
}
