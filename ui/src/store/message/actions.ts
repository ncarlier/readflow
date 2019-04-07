import { action } from 'typesafe-actions'
import { MessageActionTypes } from './types'

export const showMessage = (text: string | null, isError = false) => action(MessageActionTypes.SHOW_MESSAGE, {text, isError})
