import React, { useEffect, useRef, useState } from 'react'
import CSSTransition from 'react-transition-group/CSSTransition'

import { useMessage } from '../contexts'
import { Notification } from '.'

interface Props {
  ttl?: number
}

export const Snackbar = ({ ttl = 5000 }: Props) => {
  const { message, showMessage } = useMessage()
  const [notificationMessage, setNotificationMessage] = useState('')
  const [inProp, setInProp] = useState(false)

  const nodeRef = useRef(null)

  useEffect(() => {
    setInProp(message.text !== '')
    if (message.text !== '') {
      setNotificationMessage(message.text)
    }
  }, [message])

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
    <CSSTransition nodeRef={nodeRef} in={inProp} classNames="fade" timeout={500} unmountOnExit>
      <div ref={nodeRef}>
        <Notification message={notificationMessage} variant={message.variant}>
          <button onClick={() => showMessage('')}>dismiss</button>
        </Notification>
      </div>
    </CSSTransition>
  )
}
