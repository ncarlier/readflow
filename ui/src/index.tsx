import './index.css'

import { createBrowserHistory } from 'history'
import React from 'react'
import ReactDOM from 'react-dom'

import App from './App'
import authService from './auth'
import configureStore from './configureStore'
import { getOnlineStatus } from './helpers'
import * as serviceWorker from './serviceWorker'

const run = () => {
  const history = createBrowserHistory()
  const initialState = window.initialReduxState
  const store = configureStore(history, initialState)
  ReactDOM.render(<App store={store} history={history} />, document.getElementById('root'))
  serviceWorker.register()
  localStorage.setItem('last_run', new Date().toISOString())
}

const login = async () => {
  const user = await authService.getUser()
  if (user === null) {
    if (localStorage.getItem('last_run') === null && document.location.pathname !== '/login') {
      // No previous usage, then redirect to about page.
      document.location.replace('https://about.readflow.app')
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
  login().then(user => user && run(), () => authService.login())
} else {
  run()
}

window.addEventListener('online', login)
