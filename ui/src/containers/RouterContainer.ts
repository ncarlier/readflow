import { connect } from 'react-redux'

import { ApplicationState } from '../store'
import { RouterState, Push } from 'connected-react-router'
import { push } from 'connected-react-router'

export interface IRouterStateProps {
  router: RouterState
}

export interface IRouterDispatchProps {
  push: Push
}

const mapStateToProps = ({ router }: ApplicationState): IRouterStateProps => ({
  router
})

export const connectRouter = connect(mapStateToProps, {push})
