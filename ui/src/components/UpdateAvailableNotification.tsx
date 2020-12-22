import React from 'react'
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'

import { AppState } from '../appStore'
import { connectApp } from '../containers/AppContainer'
import Notification from './Notification'

const reload = async (registration: ServiceWorkerRegistration | null) => {
  if (registration) {
    console.log('reloading service worker...')
    await registration.update()
    document.location.reload()
  }
}

const UpdateAvailableNotification = ({ updateAvailable, registration }: AppState) => (
  <ReactCSSTransitionGroup transitionName="fade" transitionEnterTimeout={500} transitionLeaveTimeout={500}>
    {updateAvailable && (
      <Notification message="A new version is available" variant="warning">
        <button onClick={() => reload(registration)}>reload</button>
      </Notification>
    )}
  </ReactCSSTransitionGroup>
)

export default connectApp(UpdateAvailableNotification)
