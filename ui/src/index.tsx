import './index.css'

import React from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserHistory } from 'history'
import ReactModal from 'react-modal'

import App from './App'
import { isTrustedWebActivity } from './helpers'
import * as serviceWorker from './serviceWorkerRegistration'
import { PORTAL_URL } from './config'

const lastRunKey = 'readflow.lastRun'

const run = () => {
  const history = createBrowserHistory()
  ReactModal.setAppElement('#root')
  const container = document.getElementById('root')
  if (!container) {
    return
  }
  const root = createRoot(container)
  root.render(
    <React.StrictMode>
      <App history={history} />
    </React.StrictMode>,
  )
  serviceWorker.register({
    onUpdate: (registration) => {
      const event = new CustomEvent('sw-update', { detail: registration })
      window.dispatchEvent(event)
    },
  })
  localStorage.setItem(lastRunKey, new Date().toISOString())
}

const isFirstVisit = localStorage.getItem(lastRunKey) === null
const shouldRedirectToPortal =
  isFirstVisit && PORTAL_URL !== '' && !isTrustedWebActivity() && document.location.pathname !== '/login'

if (shouldRedirectToPortal) {
  document.location.replace(PORTAL_URL)
} else {
  run()
}
