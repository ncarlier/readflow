import './index.css'

import { createBrowserHistory } from 'history'
import React from 'react'
import ReactDOM from 'react-dom'

import App from './App'
import authService from './auth/AuthService'
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

authService.getUser().then(
  user => {
    if (user === null) {
      if (document.location.pathname === '/login') {
        authService.login()
      } else {
        // No previous login, then redirect to about page.
        document.location.replace('https://about.readflow.app')
      }
    } else if (user.expired) {
      authService.renewToken().then(run, () => authService.login())
    } else {
      run()
    }
  },
  () => authService.login()
)
