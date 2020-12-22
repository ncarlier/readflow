import { ApolloClient, from, HttpLink, InMemoryCache } from '@apollo/client'
import { setContext } from '@apollo/client/link/context'
import { onError } from '@apollo/client/link/error'

import authService from './auth'
import { API_BASE_URL } from './constants'

// HTTP client
const httpLink = new HttpLink({
  uri: API_BASE_URL + '/graphql',
})

// Authentication interceptor
const authenticationLink = setContext(async (_, { headers }) => {
  let user = await authService.getUser()
  if (user === null) {
    return authService.login()
  }
  if (user.expired) {
    user = await authService.renewToken()
  }
  if (user.access_token) {
    return {
      headers: {
        ...headers,
        Aauthorization: 'Bearer ' + user.access_token,
      },
    }
  }
})

// Error interceptor
const errorLink = onError((err) => {
  console.error(err)
  if (err.networkError) {
    console.log('networkError:', err.networkError.name, err.networkError.message)
    if (err.networkError.message === 'login_required') {
      authService.login()
    }
  }
})

export const client = new ApolloClient({
  link: from([errorLink, authenticationLink, httpLink]),
  cache: new InMemoryCache(),
  credentials: 'include',
})
