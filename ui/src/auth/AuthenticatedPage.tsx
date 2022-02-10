import React, { FC, useEffect } from 'react'

import { useAuth } from './AuthProvider'
import { Center, ErrorPage, Loader } from '../components'

export const AuthenticatedPage: FC = ({ children }) => {
  const { isLoading, user, error, login } = useAuth()

  useEffect(() => {
    //console.debug({ isLoading, user, error })
    if (!isLoading && !user && !error) {
      console.info('user not authenticated, redirecting to signin page...')
      login()
    }
  }, [isLoading, user, error, login])

  if (error) {
    return <ErrorPage title="Unable to authenticate user">{error?.message}</ErrorPage>
  }

  if (isLoading || !user) {
    return (
      <Center>
        <Loader />
      </Center>
    )
  }

  return <>{children}</>
}
