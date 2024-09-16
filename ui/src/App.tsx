import React from 'react'
import { ConnectedRouter } from 'connected-react-router'
import { History } from 'history'
import { ModalProvider } from 'react-modal-hook'
import { Provider } from 'react-redux'
import { Store } from 'redux'

import {
  CurrentUserProvider,
  DeviceProvider,
  GraphQLProvider,
  LocalConfigurationProvider,
  MessageProvider,
  NavbarProvider,
  ScrollMemoryProvider,
} from './contexts'
import { AppLayout } from './layout'
import Routes from './routes'
import { ApplicationState } from './store'
import { AuthenticationProvider } from './auth'

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
      <ConnectedRouter history={history}>
        <AuthenticationProvider>
          <GraphQLProvider>
            <DeviceProvider>
              <LocalConfigurationProvider>
                <MessageProvider>
                  <ModalProvider>
                    <NavbarProvider>
                      <ScrollMemoryProvider>
                        <CurrentUserProvider>
                          <AppLayout>
                            <Routes />
                          </AppLayout>
                        </CurrentUserProvider>
                      </ScrollMemoryProvider>
                    </NavbarProvider>
                  </ModalProvider>
                </MessageProvider>
              </LocalConfigurationProvider>
            </DeviceProvider>
          </GraphQLProvider>
        </AuthenticationProvider>
      </ConnectedRouter>
    </Provider>
  )
}
