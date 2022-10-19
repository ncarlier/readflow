import React, { FC, PropsWithChildren, useCallback, useMemo } from 'react'
import { ApolloClient, ApolloProvider, from, HttpLink, InMemoryCache, ServerError } from '@apollo/client'

import { API_BASE_URL } from '../constants'
import { setContext } from '@apollo/client/link/context'
import { onError } from '@apollo/client/link/error'
import { useAuth } from '../auth/AuthProvider'

// HTTP client
const httpLink = new HttpLink({
  uri: API_BASE_URL + '/graphql',
})

const cache = new InMemoryCache()

const GraphQLProvider: FC<PropsWithChildren> = ({ children }) => {
  const { user, login } = useAuth()

  // Authentication interceptor
  const authenticationLink = useCallback(() => {
    return setContext(async (_, { headers }) => {
      if (user && user.access_token) {
        return {
          headers: {
            ...headers,
            Authorization: 'Bearer ' + user.access_token,
          },
        }
      }
    })
  }, [user])

  // Error interceptor
  const errorLink = useCallback(() => {
    return onError((err) => {
      console.error(err)
      if (err.networkError) {
        const { message } = err.networkError
        if (
          (err.networkError as ServerError).statusCode === 401 ||
          message === 'login_required' ||
          message === 'invalid_grant'
        ) {
          console.warn('redirecting to login page...')
          login()
        }
      }
    })
  }, [login])

  const client = useMemo(() => {
    return new ApolloClient({
      link: from([errorLink(), authenticationLink(), httpLink]),
      cache,
      credentials: 'include',
    })
  }, [errorLink, authenticationLink])

  return <ApolloProvider client={client}>{children}</ApolloProvider>
}

export { GraphQLProvider }
