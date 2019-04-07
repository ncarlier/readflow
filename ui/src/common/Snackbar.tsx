import React, { useEffect } from 'react'
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'
import { connect } from 'react-redux'
import { Dispatch } from 'redux'

import styles from './Snackbar.module.css'
import { MessageState } from '../store/message/types'
import { ApplicationState } from '../store'
import * as messageActions from '../store/message/actions'
import { classNames } from './helpers';
import Icon from './Icon';

type Props = {
  ttl?: number
}

type AllProps = Props & IPropsFromState & IPropsFromDispatch

export const Snackbar = ({ttl = 5000, message, showMessage}: AllProps) => {
  useEffect(() => {
    if (ttl && message.text && !message.isError) {
      const timeout = setTimeout(() => {
        showMessage(null)
      }, ttl)
      return () => {
        clearTimeout(timeout);
      }
    }
  }, [ttl, message])

  const className = message.isError ? 
    classNames(styles.snackbar, styles.error) :
    styles.snackbar

  return (
    <ReactCSSTransitionGroup
      transitionName="fade"
      transitionEnterTimeout={500}
      transitionLeaveTimeout={500}>
      { message.text && 
        <div className={className}>
          <div className={styles.label}>
            {message.text}
          </div>
          <div className={styles.actions}>
            <button onClick={() => showMessage(null)}>
              dismiss
            </button>
          </div>
        </div>
      }
    </ReactCSSTransitionGroup>
  )
}

interface IPropsFromState {
  message: MessageState
}

const mapStateToProps = ({ message }: ApplicationState): IPropsFromState => ({
  message
})

interface IPropsFromDispatch {
  showMessage: typeof messageActions.showMessage
}

const mapDispatchToProps = (dispatch: Dispatch): IPropsFromDispatch => ({
  showMessage: (msg: string|null) => dispatch(messageActions.showMessage(msg))
})

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(Snackbar)
