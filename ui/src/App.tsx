import ApolloClient from 'apollo-boost'
import { ConnectedRouter } from 'connected-react-router'
import { History } from 'history'
import React from 'react'
import { ApolloProvider } from 'react-apollo-hooks'
import { ModalProvider } from 'react-modal-hook'
import { Provider } from 'react-redux'
import { Store } from 'redux'

import authService from './auth'
import { API_BASE_URL } from './constants'
import { MessageProvider } from './context/MessageContext'
import { ScrollMemoryProvider } from './context/ScrollMemoryContext'
import Routes from './routes'
import { ApplicationState } from './store'

interface PropsFromDispatch {
  [key: string]: any
}

// Any additional component props go here.
interface OwnProps {
  store: Store<ApplicationState>
  history: History
}

// Create an intersection type of the component props and our Redux props.
type Props = PropsFromDispatch & OwnProps

const client = new ApolloClient({
  uri: API_BASE_URL + '/graphql',
  fetchOptions: {
    credentials: 'include'
  },
  request: async operation => {
    let user = await authService.getUser()
    if (user === null) {
      return authService.login()
    }
    if (user.expired) {
      user = await authService.renewToken()
    }
    if (user.access_token) {
      operation.setContext({
        headers: {
          authorization: 'Bearer ' + user.access_token
        }
      })
    }
  },
  onError: err => {
    console.error(err)
    if (err.networkError) {
      console.log('networkError:', err.networkError.name, err.networkError.message)
      if (err.networkError.message === 'login_required') {
        authService.login()
      }
    }
  }
})

export default function App({ store, history /*, theme*/ }: Props) {
  return (
    <Provider store={store}>
      <ApolloProvider client={client}>
        <ModalProvider>
          <MessageProvider>
            <ConnectedRouter history={history}>
              <ScrollMemoryProvider>
                <Routes />
              </ScrollMemoryProvider>
            </ConnectedRouter>
          </MessageProvider>
        </ModalProvider>
      </ApolloProvider>
    </Provider>
  )
}
