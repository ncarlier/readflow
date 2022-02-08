import React, { FC, useEffect } from 'react'
import { AuthProvider, useAuth } from 'react-oidc-context'
import { Log } from 'oidc-client-ts'

import { AUTHORITY } from '../constants'
import { Center, ErrorPage, Loader } from '../components'
import { useOnlineStatus } from '../hooks'
import { config } from './oidc-configuration'

Log.setLogger(console)
Log.setLevel(Log.WARN)

const AuthNState: FC = ({ children }) => {
  const { isLoading, isAuthenticated, error, ...auth } = useAuth()

  useEffect(() => {
    // console.info(`AuthN ACTIVE NAVIGATOR: ${auth.activeNavigator}`)
    if (!isLoading && !isAuthenticated && !error) {
      auth.clearStaleState()
      auth.signinRedirect()
    }
  }, [isLoading, isAuthenticated, error, auth])

  switch (true) {
    case isLoading:
      return (
        <Center>
          <Loader />
        </Center>
      )
    case !!error:
      return <ErrorPage title="Unable to authenticate user">{error?.message}</ErrorPage>
    case isAuthenticated:
      return <>{children}</>
    default:
      return null
  }
}

const AuthNProvider: FC = ({ children }) => {
  const offline = !useOnlineStatus()
  const disabled = AUTHORITY === 'mock'
  if (disabled || offline) {
    return <>{children}</>
  }
  return (
    <AuthProvider {...config}>
      <AuthNState>{children}</AuthNState>
    </AuthProvider>
  )
}

export { AuthNProvider }
