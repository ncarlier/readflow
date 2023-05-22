import 'graphiql/graphiql.css'

import React, { Suspense, useCallback, useState } from 'react'
import { RouteComponentProps } from 'react-router'
import { useAuth } from '../auth/AuthProvider'

import { getAPIURL } from '../helpers'

const GraphiQL = React.lazy(() => import('graphiql'))

type AllProps = RouteComponentProps

export default ({ location }: AllProps) => {
  const query = new URLSearchParams(location.search)
  const auth = useAuth()
  const [basePath] = useState(query.has('admin') ? '/admin' : '/graphql')

  const fetcher = useCallback(
    async (graphQLParams: any) => {
      const { user } = auth
      const headers: HeadersInit = new Headers()
      headers.set('Content-Type', 'application/json')
      if (user && user.access_token) {
        headers.set('authorization', 'Bearer ' + user.access_token)
      }
      return fetch(getAPIURL(basePath), {
        method: 'post',
        headers,
        credentials: 'same-origin',
        body: JSON.stringify(graphQLParams),
      }).then((response) => response.json())
    },
    [auth, basePath]
  )

  return (
    <Suspense fallback={<div>loading...</div>}>
      <GraphiQL fetcher={fetcher} />
    </Suspense>
  )
}
