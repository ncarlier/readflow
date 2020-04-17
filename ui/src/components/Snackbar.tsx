import React, { useContext, useEffect } from 'react'
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'

import { MessageContext } from '../context/MessageContext'
import Notification from './Notification'

interface Props {
  ttl?: number
}

export default ({ ttl = 5000 }: Props) => {
  const { message, showMessage } = useContext(MessageContext)

  useEffect(() => {
    if (ttl && message.text && message.variant === 'info') {
      const timeout = setTimeout(() => {
        showMessage('')
      }, ttl)
      return () => {
        clearTimeout(timeout)
      }
    }
  }, [ttl, message, showMessage])

  return (
    <ReactCSSTransitionGroup transitionName="fade" transitionEnterTimeout={500} transitionLeaveTimeout={500}>
      {message.text && (
        <Notification message={message.text} variant={message.variant}>
          <button onClick={() => showMessage('')}>dismiss</button>
        </Notification>
      )}
    </ReactCSSTransitionGroup>
  )
}
