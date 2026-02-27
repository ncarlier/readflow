import React from 'react'
import { Router } from 'react-router-dom'
import { History } from 'history'
import { ModalProvider } from 'react-modal-hook'

import {
  CurrentUserProvider,
  DeviceProvider,
  GraphQLProvider,
  LocalConfigurationProvider,
  MessageProvider,
  NavbarProvider,
  ScrollMemoryProvider,
  AppUpdateProvider,
  OfflineArticlesProvider,
} from './contexts'
import { AppLayout } from './layout'
import Routes from './routes'
import { AuthenticationProvider } from './auth'

interface OwnProps {
  history: History
}

type Props = OwnProps

export default function App({ history /*, theme*/ }: Props) {
  return (
    <Router history={history}>
      <AuthenticationProvider>
        <GraphQLProvider>
          <DeviceProvider>
            <AppUpdateProvider>
              <LocalConfigurationProvider>
                <OfflineArticlesProvider>
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
                </OfflineArticlesProvider>
              </LocalConfigurationProvider>
            </AppUpdateProvider>
          </DeviceProvider>
        </GraphQLProvider>
      </AuthenticationProvider>
    </Router>
  )
}
