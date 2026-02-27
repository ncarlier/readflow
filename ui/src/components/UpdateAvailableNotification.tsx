import React from 'react'
import CSSTransition from 'react-transition-group/CSSTransition'

import { useAppUpdate } from '../contexts'
import { Notification } from '.'

const reload = async (registration: ServiceWorkerRegistration | null) => {
  if (registration) {
    console.log('reloading service worker...')
    await registration.update()
    document.location.reload()
  }
}

export const UpdateAvailableNotification = () => {
  const { updateAvailable, registration } = useAppUpdate()

  return (
    <CSSTransition in={updateAvailable} classNames='fade' timeout={500} unmountOnExit>
      <Notification message='A new version is available' variant='warning'>
        <button onClick={() => reload(registration)}>reload</button>
      </Notification>
    </CSSTransition>
  )
}
