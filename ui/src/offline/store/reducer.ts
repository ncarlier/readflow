import { Reducer } from 'redux'

import { Article } from '../../articles/models'
import { OfflineArticlesActionTypes, OfflineArticlesState } from './types'

// Type-safe initialState!
const initialState: OfflineArticlesState = {
  data: undefined,
  query: {
    limit: 10,
    sortBy: null,
    sortOrder: 'asc',
    afterCursor: null,
    category: null,
    starred: null,
    status: null,
    query: null,
  },
  selected: undefined,
  error: undefined,
  loading: false,
}

// Thanks to Redux 4's much simpler typings, we can take away a lot of typings on the reducer side,
// everything will remain type-safe.
const reducer: Reducer<OfflineArticlesState> = (state = initialState, action) => {
  switch (action.type) {
    case OfflineArticlesActionTypes.SAVE_REQUEST: {
      return { ...state, loading: true }
    }
    case OfflineArticlesActionTypes.SAVE_SUCCESS: {
      return { ...state, loading: false, error: undefined }
    }
    case OfflineArticlesActionTypes.SAVE_ERROR: {
      const error = action.payload
      if (error) {
        console.error(error)
      }
      return { ...state, loading: false, error }
    }
    case OfflineArticlesActionTypes.REMOVE_REQUEST: {
      return { ...state, loading: true }
    }
    case OfflineArticlesActionTypes.REMOVE_SUCCESS: {
      const article = action.payload
      let { selected } = state
      const { data } = state
      if (selected && article.id === selected.id) {
        selected = undefined
      }
      if (data) {
        data.articles.entries = data.articles.entries.filter((art: Article) => art.id !== article.id)
        data.articles.totalCount--
      }
      return { ...state, loading: false, data, error: undefined }
    }
    case OfflineArticlesActionTypes.REMOVE_ERROR: {
      return { ...state, loading: false, error: action.payload }
    }
    case OfflineArticlesActionTypes.FETCH_REQUEST: {
      const query = action.payload
      return { ...state, loading: true, query }
    }
    case OfflineArticlesActionTypes.FETCH_SUCCESS: {
      const { query, data } = state
      const { articles } = action.payload
      const nbFetchedArticles = articles.entries.length
      console.log(nbFetchedArticles + ' article(s) fetched')
      if (data && query.afterCursor) {
        articles.entries = data.articles.entries.concat(articles.entries)
      }
      return { ...state, loading: false, data: { articles }, error: undefined }
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
