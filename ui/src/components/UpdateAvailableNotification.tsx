import React from 'react'
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'

import { AppState } from '../appStore'
import { connectApp } from '../containers/AppContainer'
import Notification from './Notification'

const UpdateAvailableNotification = ({ updateAvailable }: AppState) => (
  <ReactCSSTransitionGroup transitionName="fade" transitionEnterTimeout={500} transitionLeaveTimeout={500}>
    {updateAvailable && (
      <Notification message="A new version is available" variant="warning">
        <button onClick={() => document.location.reload(true)}>reload</button>
      </Notification>
    )}
  </ReactCSSTransitionGroup>
)

export default connectApp(UpdateAvailableNotification)
