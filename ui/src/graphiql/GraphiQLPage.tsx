import 'graphiql/graphiql.css'

import React, { Suspense } from 'react'

import authService from '../auth'
import { API_BASE_URL } from '../constants'

const GraphiQL = React.lazy(() => import('graphiql'))

async function graphQLFetcher(graphQLParams: any) {
  let user = await authService.getUser()
  if (user === null) {
    return authService.login()
  }
  if (user.expired) {
    user = await authService.renewToken()
  }
  const headers: HeadersInit = new Headers()
  headers.set('Content-Type', 'application/json')
  if (user.access_token) {
    headers.set('authorization', 'Bearer ' + user.access_token)
  }
  return fetch(API_BASE_URL + '/graphql', {
    method: 'post',
    headers: headers,
    credentials: 'same-origin',
    body: JSON.stringify(graphQLParams)
  }).then(response => response.json())
}

export default () => (
  <Suspense fallback={<div>loading...</div>}>
    <GraphiQL fetcher={graphQLFetcher} />
  </Suspense>
)
