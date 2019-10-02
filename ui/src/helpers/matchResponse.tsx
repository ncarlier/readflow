import { ReactNode } from 'react'
import { ApolloError } from 'apollo-boost'

export interface GQLResponsePattern<T> {
  Loading: () => ReactNode
  Error: (err: ApolloError | Error) => ReactNode
  Data: (data: T) => ReactNode
  Other: () => ReactNode
}

export function matchResponse<T>(
  p: GQLResponsePattern<T>
): (data: T | undefined, error: ApolloError | Error | undefined, loading: boolean) => ReactNode {
  return (data: T | undefined, error: ApolloError | Error | undefined, loading: boolean): ReactNode => {
    if (loading) {
      return p.Loading()
    }
    if (error) {
      return p.Error(error)
    }
    if (data) {
      return p.Data(data)
    }
    return p.Other()
  }
}
