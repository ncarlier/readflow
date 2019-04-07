import { connect } from 'react-redux'
import { Dispatch } from 'redux'

import * as messageActions from '../store/message/actions'
import { MessageState } from '../store/message/types'
import { ApplicationState } from '../store'

export interface IMessageStateProps {
  message: MessageState
}

export interface IMessageDispatchProps {
  showMessage: typeof messageActions.showMessage
}

const mapStateToProps = ({ message }: ApplicationState): IMessageStateProps => ({
  message
})

const mapDispatchToProps = (dispatch: Dispatch): IMessageDispatchProps => ({
  showMessage: (msg: string | null) => dispatch(messageActions.showMessage(msg))
})

export const connectMessage = connect(mapStateToProps, mapDispatchToProps)
export const connectMessageDispatch = connect(null, mapDispatchToProps)
export const connectMessageState = connect(mapStateToProps)
