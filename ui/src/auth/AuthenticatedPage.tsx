import React, { FC, PropsWithChildren, useEffect } from 'react'

import { useAuth } from './AuthProvider'
import { Center, ErrorPage, Loader } from '../components'

const isLoginError = (err: any) => err && err.error && err.error === 'invalid_grant'

export const AuthenticatedPage: FC<PropsWithChildren> = ({ children }) => {
  const { isLoading, isAuthenticated, error, login } = useAuth()

  useEffect(() => {
    //console.debug({ isLoading, isAuthenticated, error })
    if (!isLoading && !error && !isAuthenticated) {
      console.warn('user not authenticated, redirecting to sign-in page...')
      login()
    }
    if (isLoginError(error)) {
      console.warn('invalid grant, redirecting to sign-in page...')
      login()
    }
  }, [isLoading, isAuthenticated, error, login])

  if (error) {
    return <ErrorPage title="Unable to authenticate user">{error?.message}</ErrorPage>
  }

  if (isLoading || !isAuthenticated) {
    return (
      <Center>
        <Loader />
      </Center>
    )
  }

  return <>{children}</>
}
