import { ReactNode } from 'react'

export interface StatePattern<T> {
  Loading: () => ReactNode
  Error: (err: Error) => ReactNode
  Data: (data: T) => ReactNode
}

export function matchState<T>(p: StatePattern<T>): (loading: boolean, data?: T, error?: Error) => ReactNode {
  return (loading: boolean, data?: T, error?: Error): ReactNode => {
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
