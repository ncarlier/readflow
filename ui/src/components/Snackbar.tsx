import React, { useContext, useEffect } from 'react'
import CSSTransition from 'react-transition-group/CSSTransition'

import { MessageContext } from '../contexts/MessageContext'
import { Notification } from '.'

interface Props {
  ttl?: number
}

export const Snackbar = ({ ttl = 5000 }: Props) => {
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
    <CSSTransition in={!!message.text} className="fade" timeout={500} unmountOnExit>
      <Notification message={message.text} variant={message.variant}>
        <button onClick={() => showMessage('')}>dismiss</button>
      </Notification>
    </CSSTransition>
  )
}
