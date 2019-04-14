import React from 'react'
import ReactDOM from 'react-dom'
import './index.css'
import App from './App'
import { createBrowserHistory } from 'history'

import configureStore from './configureStore'
import * as serviceWorker from './serviceWorker'
import authService from './auth/AuthService'

const history = createBrowserHistory()

const initialState = window.initialReduxState
const store = configureStore(history, initialState)

authService.getUser().then((user) => {
  if (user === null) {
    authService.login()
  } else {
    ReactDOM.render(<App store={store} history={history} />, document.getElementById('root'))
  }
})

serviceWorker.register()
