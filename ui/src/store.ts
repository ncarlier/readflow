import { Dispatch, Action, AnyAction } from 'redux'
import { all, fork } from 'redux-saga/effects'

import { OfflineArticlesState } from './offline/store/types'
import { offlineArticlesReducer } from './offline/store/reducer'
import { articlesSaga } from './offline/store/sagas'

import { MessageState } from './store/message/types'
import { messageReducer } from './store/message/reducer'

import { RouterState } from 'connected-react-router'

// The top-level state object.
//
// `connected-react-router` already injects the router state typings for us,
// so we can ignore them here.
export interface ApplicationState {
  offlineArticles: OfflineArticlesState
  router: RouterState
  message: MessageState
}

// Additional props for connected React components. This prop is passed by default with `connect()`
export interface ConnectedReduxProps<A extends Action = AnyAction> {
  dispatch: Dispatch<A>
}

// Whenever an action is dispatched, Redux will update each top-level application state property
// using the reducer with the matching name. It's important that the names match exactly, and that
// the reducer acts on the corresponding ApplicationState property type.
//export const rootReducer = combineReducers<ApplicationState>({
//  articles: articlesReducer
//})

export const reducers = {
  offlineArticles: offlineArticlesReducer,
  message: messageReducer
}

// Here we use `redux-saga` to trigger actions asynchronously. `redux-saga` uses something called a
// "generator function", which you can read about here:
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/function*
export function* rootSaga() {
  yield all([fork(articlesSaga)])
}
