import React, { ReactNode } from 'react'
import { ApolloError } from 'apollo-boost'

export interface GQLResponsePattern<T> {
  Loading: () => ReactNode
  Error: (err: ApolloError | Error) => ReactNode
  Data: (data: T) => ReactNode
}

export function matchResponse<T>(
  p: GQLResponsePattern<T>
): (loading: boolean, data?: T, error?: ApolloError | Error) => ReactNode {
  return (loading: boolean, data?: T, error?: ApolloError | Error): ReactNode => (
    <>
      {loading && p.Loading()}
      {error && p.Error(error)}
      {data && p.Data(data)}
    </>
  )
}
