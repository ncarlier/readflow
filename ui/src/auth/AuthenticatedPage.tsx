import React, { FC, useEffect } from 'react'

import { useAuth } from './AuthProvider'
import { Center, ErrorPage, Loader } from '../components'

export const AuthenticatedPage: FC = ({ children }) => {
  const { isLoading, isAuthenticated, error, login } = useAuth()

  useEffect(() => {
    console.debug({ isLoading, isAuthenticated, error })
    if (!isLoading && !error && !isAuthenticated) {
      console.info('user not authenticated, redirecting to sign-in page...')
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
