import React, { createContext, FC, PropsWithChildren, useContext, useCallback, useReducer } from 'react'
import { Article, GetArticlesRequest, GetArticlesResponse } from '../articles/models'
import { getArticle, getArticles, removeArticle, saveArticle } from '../offline/dao/articles'

export interface OfflineArticlesState {
  loading: boolean
  data?: GetArticlesResponse
  query: GetArticlesRequest
  selected?: Article
  error?: Error
}

type OfflineArticlesAction =
  | { type: 'FETCH_REQUEST'; payload: GetArticlesRequest }
  | { type: 'FETCH_SUCCESS'; payload: GetArticlesResponse }
  | { type: 'FETCH_ERROR'; payload: Error }
  | { type: 'SELECT_REQUEST'; payload: number }
  | { type: 'SELECT_SUCCESS'; payload: Article | undefined }
  | { type: 'SELECT_ERROR'; payload: Error }
  | { type: 'SAVE_REQUEST'; payload: Article }
  | { type: 'SAVE_SUCCESS'; payload: Article }
  | { type: 'SAVE_ERROR'; payload: Error }
  | { type: 'REMOVE_REQUEST'; payload: Article }
  | { type: 'REMOVE_SUCCESS'; payload: Article }
  | { type: 'REMOVE_ERROR'; payload: Error }

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

function reducer(state: OfflineArticlesState, action: OfflineArticlesAction): OfflineArticlesState {
  switch (action.type) {
    case 'SAVE_REQUEST':
      return { ...state, loading: true }
    case 'SAVE_SUCCESS':
      return { ...state, loading: false, error: undefined }
    case 'SAVE_ERROR':
      return { ...state, loading: false, error: action.payload }
    case 'REMOVE_REQUEST':
      return { ...state, loading: true }
    case 'REMOVE_SUCCESS': {
      const article = action.payload
      let { selected } = state
      const data = state.data ? { ...state.data } : undefined
      if (selected && article.id === selected.id) {
        selected = undefined
      }
      if (data) {
        data.articles.entries = data.articles.entries.filter((art: Article) => art.id !== article.id)
        data.articles.totalCount--
      }
      return { ...state, loading: false, data, selected, error: undefined }
    }
    case 'REMOVE_ERROR':
      return { ...state, loading: false, error: action.payload }
    case 'FETCH_REQUEST':
      return { ...state, loading: true, query: action.payload }
    case 'FETCH_SUCCESS': {
      const { query, data } = state
      const { articles } = action.payload
      if (data && query.afterCursor) {
        articles.entries = data.articles.entries.concat(articles.entries)
      }
      return { ...state, loading: false, data: { articles }, error: undefined }
    }
    case 'FETCH_ERROR':
      return { ...state, loading: false, error: action.payload }
    case 'SELECT_REQUEST':
      return { ...state, loading: true, selected: undefined }
    case 'SELECT_SUCCESS':
      return { ...state, loading: false, selected: action.payload }
    case 'SELECT_ERROR':
      return { ...state, loading: false, error: action.payload }
    default:
      return state
  }
}

interface OfflineArticlesContextType {
  offlineArticles: OfflineArticlesState
  fetchOfflineArticles: (query: GetArticlesRequest) => Promise<void>
  fetchOfflineArticle: (id: number) => Promise<void>
  saveOfflineArticle: (article: Article) => Promise<void>
  removeOfflineArticle: (article: Article) => Promise<void>
}

const OfflineArticlesContext = createContext<OfflineArticlesContextType>({
  offlineArticles: initialState,
  fetchOfflineArticles: () => Promise.resolve(),
  fetchOfflineArticle: () => Promise.resolve(),
  saveOfflineArticle: () => Promise.resolve(),
  removeOfflineArticle: () => Promise.resolve(),
})

export const OfflineArticlesProvider: FC<PropsWithChildren> = ({ children }) => {
  const [offlineArticles, dispatch] = useReducer(reducer, initialState)

  const fetchOfflineArticles = useCallback(async (query: GetArticlesRequest) => {
    dispatch({ type: 'FETCH_REQUEST', payload: query })
    try {
      const res = await getArticles(query)
      dispatch({ type: 'FETCH_SUCCESS', payload: res })
    } catch (err) {
      dispatch({ type: 'FETCH_ERROR', payload: err instanceof Error ? err : new Error('Unknown error') })
    }
  }, [])

  const fetchOfflineArticle = useCallback(async (id: number) => {
    dispatch({ type: 'SELECT_REQUEST', payload: id })
    try {
      const res = await getArticle(id)
      dispatch({ type: 'SELECT_SUCCESS', payload: res })
    } catch (err) {
      dispatch({ type: 'SELECT_ERROR', payload: err instanceof Error ? err : new Error('Unknown error') })
    }
  }, [])

  const saveOfflineArticle = useCallback(async (article: Article) => {
    dispatch({ type: 'SAVE_REQUEST', payload: article })
    try {
      const res = await saveArticle({ ...article, isOffline: true })
      dispatch({ type: 'SAVE_SUCCESS', payload: res })
    } catch (err) {
      dispatch({ type: 'SAVE_ERROR', payload: err instanceof Error ? err : new Error('Unknown error') })
      throw err // Rethrow to allow caller to handle
    }
  }, [])

  const removeOfflineArticle = useCallback(async (article: Article) => {
    dispatch({ type: 'REMOVE_REQUEST', payload: article })
    try {
      const res = await removeArticle(article)
      dispatch({ type: 'REMOVE_SUCCESS', payload: res })
    } catch (err) {
      dispatch({ type: 'REMOVE_ERROR', payload: err instanceof Error ? err : new Error('Unknown error') })
      throw err // Rethrow to allow caller to handle
    }
  }, [])

  return (
    <OfflineArticlesContext.Provider
      value={{ offlineArticles, fetchOfflineArticles, fetchOfflineArticle, saveOfflineArticle, removeOfflineArticle }}
    >
      {children}
    </OfflineArticlesContext.Provider>
  )
}

export const useOfflineArticles = () => useContext(OfflineArticlesContext)
