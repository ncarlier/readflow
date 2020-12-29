import { ApolloError } from '@apollo/client'
import { ReactNode } from 'react'

export interface GQLResponsePattern<T> {
  Loading: () => ReactNode
  Error: (err: ApolloError | Error) => ReactNode
  Data: (data: T) => ReactNode
}

export function matchResponse<T>(
  p: GQLResponsePattern<T>
): (loading: boolean, data?: T, error?: ApolloError | Error) => ReactNode {
  return (loading: boolean, data?: T, error?: ApolloError | Error): ReactNode => {
    if (loading) {
      return p.Loading()
    }
    if (error !== undefined) {
      return p.Error(error)
    }
    if (data !== undefined) {
      return p.Data(data)
    }
    return null
  }
}
