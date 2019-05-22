import React, { useContext, useEffect } from 'react'
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'

import { MessageContext } from '../context/MessageContext'
import { classNames } from './helpers'
import styles from './Snackbar.module.css'

interface Props {
  ttl?: number
}

export default ({ ttl = 5000 }: Props) => {
  const { message, showMessage } = useContext(MessageContext)

  useEffect(() => {
    if (ttl && message.text && !message.isError) {
      const timeout = setTimeout(() => {
        showMessage('')
      }, ttl)
      return () => {
        clearTimeout(timeout)
      }
    }
  }, [ttl, message])

  const className = message.isError ? classNames(styles.snackbar, styles.error) : styles.snackbar

  return (
    <ReactCSSTransitionGroup transitionName="fade" transitionEnterTimeout={500} transitionLeaveTimeout={500}>
      {message.text && (
        <div className={className}>
          <div className={styles.label}>{message.text}</div>
          <div className={styles.actions}>
            <button onClick={() => showMessage('')}>dismiss</button>
          </div>
        </div>
      )}
    </ReactCSSTransitionGroup>
  )
}
