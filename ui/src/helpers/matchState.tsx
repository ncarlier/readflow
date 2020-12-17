import React, { ReactNode } from 'react'

export interface StatePattern<T> {
  Loading: () => ReactNode
  Error: (err: Error) => ReactNode
  Data: (data: T) => ReactNode
}

export function matchState<T>(p: StatePattern<T>): (loading: boolean, data?: T, error?: Error) => ReactNode {
  return (loading: boolean, data?: T, error?: Error): ReactNode => (
    <>
      {loading && p.Loading()}
      {error && p.Error(error)}
      {data && p.Data(data)}
    </>
  )
}
