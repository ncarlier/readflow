import './index.css'

import { createBrowserHistory } from 'history'
import React from 'react'
import ReactDOM from 'react-dom'

import App from './App'
import authService from './auth/AuthService'
import { getOnlineStatus } from './common/helpers'
import configureStore from './configureStore'
import { setupNotification } from './notification'
import * as serviceWorker from './serviceWorker'

const run = () => {
  const history = createBrowserHistory()
  const initialState = window.initialReduxState
  const store = configureStore(history, initialState)
  ReactDOM.render(<App store={store} history={history} />, document.getElementById('root'))
  serviceWorker.register()
  setupNotification()
}

const login = async () => {
  const user = await authService.getUser()
  if (user === null) {
    if (document.location.pathname === '/login') {
      return authService.login()
    } else {
      // No previous login, then redirect to about page.
      document.location.replace('https://about.readflow.app')
    }
  } else if (user.expired) {
    return authService.renewToken()
  } else {
    return user
  }
}

if (getOnlineStatus()) {
  login().then(run, authService.login)
} else {
  run()
}

window.addEventListener('online', login)
