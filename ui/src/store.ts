import { RouterState } from 'connected-react-router'
import { Action, AnyAction, Dispatch } from 'redux'
import { all, fork } from 'redux-saga/effects'

import { offlineArticlesReducer } from './offline/store/reducer'
import { articlesSaga } from './offline/store/sagas'
import { OfflineArticlesState } from './offline/store/types'
import { AppState, appReducer } from './appStore'

// The top-level state object.
//
// `connected-react-router` already injects the router state typings for us,
// so we can ignore them here.
export interface ApplicationState {
  app: AppState
  offlineArticles: OfflineArticlesState
  router: RouterState
}

// Additional props for connected React components. This prop is passed by default with `connect()`
export interface ConnectedReduxProps<A extends Action = AnyAction> {
  dispatch: Dispatch<A>
}

export const reducers = {
  app: appReducer,
  offlineArticles: offlineArticlesReducer
}

// Here we use `redux-saga` to trigger actions asynchronously. `redux-saga` uses something called a
// "generator function", which you can read about here:
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/function*
export function* rootSaga() {
  yield all([fork(articlesSaga)])
}
