import './index.css'

import React from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserHistory } from 'history'
import ReactModal from 'react-modal'

import App from './App'
import { updateAvailable } from './appStore'
import configureStore from './configureStore'
import { isTrustedWebActivity } from './helpers'
import * as serviceWorker from './serviceWorkerRegistration'
import { REDIRECT_URL } from './constants'
import { ApplicationState } from './store'

const lastRunKey = 'readflow.lastRun'

const run = () => {
  const history = createBrowserHistory()
  const initialState = window.initialReduxState as ApplicationState
  const store = configureStore(history, initialState)
  ReactModal.setAppElement('#root')
  const container = document.getElementById('root')
  if (!container) {
    return
  }
  const root = createRoot(container)
  root.render(
    <React.StrictMode>
      <App store={store} history={history} />
    </React.StrictMode>
  )
  serviceWorker.register({ onUpdate: (registration) => store.dispatch(updateAvailable(registration)) })
  localStorage.setItem(lastRunKey, new Date().toISOString())
}

const shouldRedirectToPortal = () =>
  !isTrustedWebActivity() && localStorage.getItem(lastRunKey) === null && document.location.pathname !== '/login'

if (shouldRedirectToPortal()) {
  document.location.replace(REDIRECT_URL)
} else {
  run()
}
