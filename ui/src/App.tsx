import React from 'react'
import { ApolloProvider } from '@apollo/client'
import { ConnectedRouter } from 'connected-react-router'
import { History } from 'history'
import { ModalProvider } from 'react-modal-hook'
import { Provider } from 'react-redux'
import { Store } from 'redux'

import {
  CurrentUserProvider,
  LocalConfigurationProvider,
  MessageProvider,
  NavbarProvider,
  ScrollMemoryProvider,
} from './contexts'
import { client } from './graphqlClient'
import { AppLayout } from './layout'
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

export default function App({ store, history /*, theme*/ }: Props) {
  return (
    <Provider store={store}>
      <ApolloProvider client={client}>
        <LocalConfigurationProvider>
          <ModalProvider>
            <MessageProvider>
              <NavbarProvider>
                <ScrollMemoryProvider>
                  <ConnectedRouter history={history}>
                    <CurrentUserProvider>
                      <AppLayout>
                        <Routes />
                      </AppLayout>
                    </CurrentUserProvider>
                  </ConnectedRouter>
                </ScrollMemoryProvider>
              </NavbarProvider>
            </MessageProvider>
          </ModalProvider>
        </LocalConfigurationProvider>
      </ApolloProvider>
    </Provider>
  )
}
