import React, { FC, PropsWithChildren } from 'react'
import { AuthProvider } from './AuthProvider'

import { AuthenticatedPage } from './AuthenticatedPage'
import { AUTHORITY } from '../config'
import { useOnlineStatus } from '../hooks'

const AuthenticationProvider: FC<PropsWithChildren> = ({ children }) => {
  const offline = !useOnlineStatus()
  const disabled = AUTHORITY === 'none'
  if (disabled || offline) {
    return <>{children}</>
  }
  return (
    <AuthProvider>
      <AuthenticatedPage>{children}</AuthenticatedPage>
    </AuthProvider>
  )
}

export { AuthenticationProvider }
