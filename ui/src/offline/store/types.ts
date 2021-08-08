import { Article, GetArticlesRequest, GetArticlesResponse } from '../../articles/models'

export enum OfflineArticlesActionTypes {
  SAVE_REQUEST = '@@offlineArticles/SAVE_REQUEST',
  SAVE_SUCCESS = '@@offlineArticles/SAVE_SUCCESS',
  SAVE_ERROR = '@@offlineArticles/SAVE_ERROR',
  REMOVE_REQUEST = '@@offlineArticles/REMOVE_REQUEST',
  REMOVE_SUCCESS = '@@offlineArticles/REMOVE_SUCCESS',
  REMOVE_ERROR = '@@offlineArticles/REMOVE_ERROR',
  FETCH_REQUEST = '@@offlineArticles/FETCH_REQUEST',
  FETCH_SUCCESS = '@@offlinAarticles/FETCH_SUCCESS',
  FETCH_ERROR = '@@offlineArticles/FETCH_ERROR',
  SELECT_REQUEST = '@@offlineArticles/SELECT_REQUEST',
  SELECT_SUCCESS = '@@offlinAarticles/SELECT_SUCCESS',
  SELECT_ERROR = '@@offlineArticles/SELECT_ERROR',
}

// Declare state types with `readonly` modifier to get compile time immutability.
// https://github.com/piotrwitek/react-redux-typescript-guide#state-with-type-level-immutability
export interface OfflineArticlesState {
  readonly loading: boolean
  readonly data?: GetArticlesResponse
  readonly query: GetArticlesRequest
  readonly selected?: Article
  readonly error?: Error
}

export interface ErrorResponse {
  readonly error: any
}
