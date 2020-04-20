import { Reducer } from 'redux'
import { action } from 'typesafe-actions'

export enum AppActionTypes {
  UPDATE_AVAILABLE = '@@app/UPDATE_AVAILABLE',
}

export interface AppState {
  readonly updateAvailable: boolean
  readonly registration: ServiceWorkerRegistration | null
}

export const updateAvailable = (registration: ServiceWorkerRegistration) =>
  action(AppActionTypes.UPDATE_AVAILABLE, registration)

const initialState: AppState = {
  updateAvailable: false,
  registration: null,
}

const reducer: Reducer<AppState> = (state = initialState, action) => {
  switch (action.type) {
    case AppActionTypes.UPDATE_AVAILABLE: {
      const registration = action.payload
      return { ...state, updateAvailable: true, registration }
    }
    default: {
      return state
    }
  }
}

export { reducer as appReducer }
