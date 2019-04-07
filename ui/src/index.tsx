import React from 'react'
import ReactDOM from 'react-dom'
import './index.css'
import App from './App'
import { createBrowserHistory } from 'history'

import configureStore from './configureStore'
import * as serviceWorker from './serviceWorker'


// We use hash history because this example is going to be hosted statically.
// Normally you would use browser history.
const history = createBrowserHistory()

const initialState = window.initialReduxState
const store = configureStore(history, initialState)

ReactDOM.render(<App store={store} history={history} />, document.getElementById('root'))

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister()
