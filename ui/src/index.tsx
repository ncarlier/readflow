import './index.css'

import { createBrowserHistory } from 'history'
import React from 'react'
import ReactDOM from 'react-dom'

import App from './App'
import { updateAvailable } from './appStore'
import authService from './auth'
import configureStore from './configureStore'
import { getOnlineStatus, isTrustedWebActivity } from './helpers'
import * as serviceWorker from './serviceWorker'
import { REDIRECT_URL } from './constants'

const lastRunKey = 'readflow.lastRun'

const run = () => {
  const history = createBrowserHistory()
  const initialState = window.initialReduxState
  const store = configureStore(history, initialState)
  ReactDOM.render(<App store={store} history={history} />, document.getElementById('root'))
  serviceWorker.register({ onUpdate: (registration) => store.dispatch(updateAvailable(registration)) })
  localStorage.setItem(lastRunKey, new Date().toISOString())
}

const shouldRedirect = () => {
  return !isTrustedWebActivity() && localStorage.getItem(lastRunKey) === null && document.location.pathname !== '/login'
}

const login = async () => {
  const user = await authService.getUser()
  if (user === null) {
    if (shouldRedirect()) {
      // No previous usage, then redirect to about page.
      document.location.replace(REDIRECT_URL)
    } else {
      throw new Error('login forced')
    }
  } else if (user.expired) {
    return await authService.renewToken()
  } else {
    return Promise.resolve(user)
  }
}

if (getOnlineStatus()) {
  login().then(
    (user) => user && run(),
    () => authService.login()
  )
} else {
  run()
}

window.addEventListener('online', login)
