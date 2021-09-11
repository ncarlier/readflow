import React, { createContext, FC, useContext } from 'react'
import { gql, useQuery } from '@apollo/client'
import { getGQLError, matchResponse } from '../helpers'
import { Center, ErrorPage, Loader } from '../components'
import { useOnlineStatus } from '../hooks'

const GetCurrentUser = gql`
  query {
    me {
      username
      hash
      hashid
      plan
      customer_id
      last_login_at
      created_at
    }
  }
`

export interface User {
  username: string
  hash: string
  hashid: string
  plan: string
  customer_id: string
  created_at: string
  last_login_at: string
}

interface GetCurrentUserResponse {
  me: User
}

const CurrentUserContext = createContext<User | null>(null)

const CurrentUserProvider: FC = ({ children }) => {
  const online = useOnlineStatus()
  const { data, error, loading } = useQuery<GetCurrentUserResponse>(GetCurrentUser)

  const render = matchResponse<GetCurrentUserResponse>({
    Loading: () => (
      <Center>
        <Loader />
      </Center>
    ),
    Error: (err) => <ErrorPage title="Unable to get user information">{getGQLError(err)}</ErrorPage>,
    Data: ({ me: user }) => <CurrentUserContext.Provider value={user}>{children}</CurrentUserContext.Provider>,
  })

  if (!online) {
    return <CurrentUserContext.Provider value={null}>{children}</CurrentUserContext.Provider>
  }

  return <>{render(loading, data, error)}</>
}

export { CurrentUserProvider }

export const useCurrentUser = () => useContext(CurrentUserContext)
