import React, { FC, useContext } from 'react'
import { AuthContext } from 'react-oidc-context'
import { ApolloClient, ApolloProvider, from, HttpLink, InMemoryCache } from '@apollo/client'

import { API_BASE_URL } from '../constants'
import { setContext } from '@apollo/client/link/context'
import { onError } from '@apollo/client/link/error'

// HTTP client
const httpLink = new HttpLink({
  uri: API_BASE_URL + '/graphql',
})

const GraphQLProvider: FC = ({ children }) => {
  const auth = useContext(AuthContext)

  // Authentication interceptor
  const authenticationLink = setContext(async (_, { headers }) => {
    const user = auth?.user
    if (user && user.access_token) {
      return {
        headers: {
          ...headers,
          Authorization: 'Bearer ' + user.access_token,
        },
      }
    }
  })
  // Error interceptor
  const errorLink = onError((err) => {
    console.error(err)
    if (err.networkError) {
      const { name, message } = err.networkError
      console.error('networkError:', name, message)
      if (message === 'login_required' || message === 'invalid_grant') {
        console.warn('redirecting to login page...')
        auth?.signinRedirect()
      }
    }
  })

  const client = new ApolloClient({
    link: from([errorLink, authenticationLink, httpLink]),
    cache: new InMemoryCache(),
    credentials: 'include',
  })
  return <ApolloProvider client={client}>{children}</ApolloProvider>
}

export { GraphQLProvider }
