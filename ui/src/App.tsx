import React from 'react'
import { Provider} from 'react-redux'
import { ModalProvider } from "react-modal-hook"
import { ConnectedRouter } from 'connected-react-router'
import { Store } from 'redux'
import { History } from 'history'

import ApolloClient from "apollo-boost"
import { ApolloProvider } from "react-apollo-hooks"

import Routes from './routes'
import { ApplicationState } from './store'
import authService from './auth/AuthService'
import { API_BASE_URL } from './constants'

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
  request: async (operation) => {
    let user = await authService.getUser()
    if (user === null) {
      authService.login()
    }
    if (user.expired) {
      user = await authService.renewToken()
    }
    operation.setContext({
      headers: {
        authorization: 'Bearer ' + user.access_token
      }
    })
  },
  onError: ({ graphQLErrors, networkError, response }) => {
    console.log('onError', graphQLErrors, networkError, response)
  },
})

export default function App({store, history, theme}: Props) {
  return (
    <Provider store={store}>
      <ApolloProvider client={client}>
        <ModalProvider>
          <ConnectedRouter history={history}>
            <Routes />
          </ConnectedRouter>
        </ModalProvider>
      </ApolloProvider>
    </Provider>
  )
}
