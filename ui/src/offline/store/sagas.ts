import { all, call, fork, put, takeEvery } from 'redux-saga/effects'

import { getArticle, getArticles, removeArticle, saveArticle } from '../dao/articles'
import {
  fetchError,
  fetchRequest,
  fetchSuccess,
  removeError,
  removeRequest,
  removeSuccess,
  saveError,
  saveRequest,
  saveSuccess,
  selectError,
  selectRequest,
  selectSuccess,
} from './actions'
import { OfflineArticlesActionTypes } from './types'

// Here we use `redux-saga` to trigger actions asynchronously. `redux-saga` uses something called a
// "generator function", which you can read about here:
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/function*

function* handleFetch({ payload }: ReturnType<typeof fetchRequest>) {
  try {
    // To call async functions, use redux-saga's `call()`.
    const res = yield call(getArticles, payload)

    if (res.error) {
      yield put(fetchError(new Error(res.error)))
    } else {
      yield put(fetchSuccess(res))
    }
  } catch (err) {
    if (err instanceof Error) {
      yield put(fetchError(err))
    } else {
      yield put(fetchError(new Error('An unknown error occured.')))
    }
  }
}

function* handleSave({ payload }: ReturnType<typeof saveRequest>) {
  try {
    const res = yield call(saveArticle, { ...payload, isOffline: true })

    if (res.error) {
      yield put(saveError(new Error(res.error)))
    } else {
      yield put(saveSuccess(res))
    }
  } catch (err) {
    if (err instanceof Error) {
      yield put(saveError(err))
    } else {
      yield put(saveError(new Error('An unknown error occured.')))
    }
  }
}

function* handleRemove({ payload }: ReturnType<typeof removeRequest>) {
  try {
    const res = yield call(removeArticle, payload)

    if (res.error) {
      yield put(removeError(new Error(res.error)))
    } else {
      yield put(removeSuccess(res))
    }
  } catch (err) {
    if (err instanceof Error) {
      yield put(removeError(err))
    } else {
      yield put(removeError(new Error('An unknown error occured.')))
    }
  }
}

function* handleSelect({ payload }: ReturnType<typeof selectRequest>) {
  try {
    // To call async functions, use redux-saga's `call()`.
    const res = yield call(getArticle, payload)

    if (res.error) {
      yield put(selectError(new Error(res.error)))
    } else {
      yield put(selectSuccess(res))
    }
  } catch (err) {
    if (err instanceof Error) {
      yield put(selectError(err))
    } else {
      yield put(selectError(new Error('An unknown error occured.')))
    }
  }
}

// This is our watcher function. We use `take*()` functions to watch Redux for a specific action
// type, and run our saga, for example the `handleFetch()` saga above.
function* watchFetchRequest() {
  yield takeEvery(OfflineArticlesActionTypes.FETCH_REQUEST, handleFetch)
}

function* watchSelectRequest() {
  yield takeEvery(OfflineArticlesActionTypes.SELECT_REQUEST, handleSelect)
}

function* watchSaveRequest() {
  yield takeEvery(OfflineArticlesActionTypes.SAVE_REQUEST, handleSave)
}

function* watchRemoveRequest() {
  yield takeEvery(OfflineArticlesActionTypes.REMOVE_REQUEST, handleRemove)
}

// Export our root saga.
// We can also use `fork()` here to split our saga into multiple watchers.
export function* articlesSaga() {
  yield all([fork(watchFetchRequest), fork(watchSelectRequest), fork(watchSaveRequest), fork(watchRemoveRequest)])
}
