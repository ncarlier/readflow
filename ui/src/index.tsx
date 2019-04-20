import React from 'react'
import ReactDOM from 'react-dom'
import './index.css'
import App from './App'
import { createBrowserHistory } from 'history'

import configureStore from './configureStore'
import * as serviceWorker from './serviceWorker'
import authService from './auth/AuthService'

const run = () => {
  const history = createBrowserHistory()
  const initialState = window.initialReduxState
  const store = configureStore(history, initialState)
  ReactDOM.render(<App store={store} history={history} />, document.getElementById('root'))
  serviceWorker.register()
}

authService.getUser().then((user) => {
  if (user === null) {
    if (document.location.pathname === '/login') {
      authService.login()
    } else {
      // No previous login, then redirect to about page.
      document.location.replace("https://about.readflow.app")
    }
  } else if (user.expired) {
    authService.renewToken().then(run)
  } else {
    run()
  }
})
