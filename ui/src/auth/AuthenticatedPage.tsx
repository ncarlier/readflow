import React, { FC, PropsWithChildren } from 'react'

import { useAuth } from './AuthProvider'
import { Center, ErrorPage, Loader } from '../components'

//const isLoginError = (err: any) => err && err.error && err.error === 'invalid_grant'

export const AuthenticatedPage: FC<PropsWithChildren> = ({ children }) => {
  const { isLoading, user, error } = useAuth()

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
