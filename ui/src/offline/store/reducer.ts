import { Reducer } from 'redux'
import { OfflineArticlesState, OfflineArticlesActionTypes } from './types'
import { Article } from '../../articles/models'

// Type-safe initialState!
const initialState: OfflineArticlesState = {
  data: [],
  selected: undefined,
  error: undefined,
  loading: true
}

// Thanks to Redux 4's much simpler typings, we can take away a lot of typings on the reducer side,
// everything will remain type-safe.
const reducer: Reducer<OfflineArticlesState> = (state = initialState, action) => {
  switch (action.type) {
    case OfflineArticlesActionTypes.SAVE_REQUEST: {
      return { ...state, loading: true }
    }
    case OfflineArticlesActionTypes.SAVE_SUCCESS: {
      let { data } = state
      const article = action.payload
      data = data.filter((art: Article) => art.id != article.id)
      data.unshift(action.payload)
      return { ...state, loading: false, data }
    }
    case OfflineArticlesActionTypes.SAVE_ERROR: {
      return { ...state, loading: false, error: action.payload }
    }
    case OfflineArticlesActionTypes.REMOVE_REQUEST: {
      return { ...state, loading: true }
    }
    case OfflineArticlesActionTypes.REMOVE_SUCCESS: {
      const article = action.payload
      let { data, selected } = state
      if (selected && article.id === selected.id) {
        selected = undefined
      }
      data = data.filter((art: Article) => art.id != article.id)
      return { ...state, loading: false, data, selected }
    }
    case OfflineArticlesActionTypes.REMOVE_ERROR: {
      return { ...state, loading: false, error: action.payload }
    }
    case OfflineArticlesActionTypes.FETCH_REQUEST: {
      return { ...state, loading: true }
    }
    case OfflineArticlesActionTypes.FETCH_SUCCESS: {
      return { ...state, loading: false, data: action.payload }
    }
    case OfflineArticlesActionTypes.FETCH_ERROR: {
      return { ...state, loading: false, error: action.payload }
    }
    case OfflineArticlesActionTypes.SELECT_REQUEST: {
      return { ...state, loading: true, selected: undefined }
    }
    case OfflineArticlesActionTypes.SELECT_SUCCESS: {
      return { ...state, loading: false, selected: action.payload }
    }
    case OfflineArticlesActionTypes.SELECT_ERROR: {
      return { ...state, loading: false, error: action.payload }
    }
    default: {
      return state
    }
  }
}

// Instead of using default export, we use named exports. That way we can group these exports
// inside the `index.js` folder.
export { reducer as offlineArticlesReducer }
