import React, { FC, useEffect } from 'react'
import { AuthProvider, useAuth } from 'react-oidc-context'

import { AUTHORITY } from '../constants'
import { Center, ErrorPage, Loader } from '../components'
import { useOnlineStatus } from '../hooks'
import { buildConfig } from './oidc-configuration'

const AuthNState: FC = ({ children }) => {
  const { isLoading, isAuthenticated, error, ...auth } = useAuth()

  useEffect(() => {
    // console.debug(`authn active navigator: ${auth.activeNavigator}`)
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
  const redirect = encodeURIComponent(document.location.href)
  return (
    <AuthProvider {...buildConfig(redirect)}>
      <AuthNState>{children}</AuthNState>
    </AuthProvider>
  )
}

export { AuthNProvider }
