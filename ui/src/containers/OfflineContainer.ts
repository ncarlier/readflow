import { connect } from 'react-redux'
import { Dispatch } from 'redux'

import { Article } from '../articles/models'
import { GetArticlesQuery } from '../offline/dao/articles'
import * as offlineArticlesActions from '../offline/store/actions'
import { OfflineArticlesState } from '../offline/store/types'
import { ApplicationState } from '../store'

export interface IOfflineStateProps {
  offlineArticles: OfflineArticlesState
}

export interface IOfflineDispatchProps {
  saveOfflineArticle: typeof offlineArticlesActions.saveRequest
  removeOfflineArticle: typeof offlineArticlesActions.removeRequest
  fetchOfflineArticles: typeof offlineArticlesActions.fetchRequest
  fetchOfflineArticle: typeof offlineArticlesActions.selectRequest
}

export type OfflineProps = IOfflineStateProps & IOfflineDispatchProps

const mapStateToProps = ({ offlineArticles }: ApplicationState): IOfflineStateProps => ({
  offlineArticles
})

const mapDispatchToProps = (dispatch: Dispatch): IOfflineDispatchProps => ({
  saveOfflineArticle: (data: Article) => dispatch(offlineArticlesActions.saveRequest(data)),
  removeOfflineArticle: (data: Article) => dispatch(offlineArticlesActions.removeRequest(data)),
  fetchOfflineArticles: (query: GetArticlesQuery) => dispatch(offlineArticlesActions.fetchRequest(query)),
  fetchOfflineArticle: (id: number) => dispatch(offlineArticlesActions.selectRequest(id))
})

export const connectOffline = connect(
  mapStateToProps,
  mapDispatchToProps
)
export const connectOfflineDispatch = connect(
  null,
  mapDispatchToProps
)
export const connectOfflineState = connect(mapStateToProps)
