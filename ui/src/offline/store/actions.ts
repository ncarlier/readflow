import { action } from 'typesafe-actions'

import { Article, GetArticlesRequest, GetArticlesResponse } from '../../articles/models'
import { OfflineArticlesActionTypes } from './types'

export const saveRequest = (data: Article) => action(OfflineArticlesActionTypes.SAVE_REQUEST, data)
export const saveSuccess = (data: Article) => action(OfflineArticlesActionTypes.SAVE_SUCCESS, data)
export const saveError = (err: Error) => action(OfflineArticlesActionTypes.SAVE_ERROR, err)

export const removeRequest = (data: Article) => action(OfflineArticlesActionTypes.REMOVE_REQUEST, data)
export const removeSuccess = (data: Article) => action(OfflineArticlesActionTypes.REMOVE_SUCCESS, data)
export const removeError = (err: Error) => action(OfflineArticlesActionTypes.REMOVE_ERROR, err)

export const fetchRequest = (query: GetArticlesRequest) => action(OfflineArticlesActionTypes.FETCH_REQUEST, query)
export const fetchSuccess = (data: GetArticlesResponse) => action(OfflineArticlesActionTypes.FETCH_SUCCESS, data)
export const fetchError = (err: Error) => action(OfflineArticlesActionTypes.FETCH_ERROR, err)

export const selectRequest = (id: number) => action(OfflineArticlesActionTypes.SELECT_REQUEST, id)
export const selectSuccess = (data: Article) => action(OfflineArticlesActionTypes.SELECT_SUCCESS, data)
export const selectError = (err: Error) => action(OfflineArticlesActionTypes.SELECT_ERROR, err)
