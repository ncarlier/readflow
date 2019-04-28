import { Push, push, RouterState } from 'connected-react-router'
import { connect } from 'react-redux'

import { ApplicationState } from '../store'

export interface IRouterStateProps {
  router: RouterState
}

export interface IRouterDispatchProps {
  push: Push
}

const mapStateToProps = ({ router }: ApplicationState): IRouterStateProps => ({
  router
})

export const connectRouter = connect(
  mapStateToProps,
  { push }
)
