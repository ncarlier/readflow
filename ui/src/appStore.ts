import { action } from 'typesafe-actions'
import { Reducer } from 'redux'

export enum AppActionTypes {
  UPDATE_AVAILABLE = '@@app/UPDATE_AVAILABLE'
}

export interface AppState {
  readonly updateAvailable: boolean
}

export const updateAvailable = () => action(AppActionTypes.UPDATE_AVAILABLE)

const initialState: AppState = {
  updateAvailable: false
}

const reducer: Reducer<AppState> = (state = initialState, action) => {
  switch (action.type) {
    case AppActionTypes.UPDATE_AVAILABLE: {
      return { ...state, updateAvailable: true }
    }
    default: {
      return state
    }
  }
}

export { reducer as appReducer }
