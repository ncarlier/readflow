import { Reducer } from 'redux'

import { MessageActionTypes, MessageState } from './types'

const initialState: MessageState = {
  text: null,
  isError: false
}

const reducer: Reducer<MessageState> = (state = initialState, action) => {
  switch (action.type) {
    case MessageActionTypes.SHOW_MESSAGE: {
      return { ...state, ...action.payload }
    }
    default: {
      return state
    }
  }
}

export { reducer as messageReducer }
