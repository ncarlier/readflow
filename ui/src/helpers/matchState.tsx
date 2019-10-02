import React, { ReactNode } from 'react'

export interface StatePattern<T> {
  Loading: () => ReactNode
  Error: (err: Error) => ReactNode
  Data: (data: T) => ReactNode
}

export function matchState<T>(
  p: StatePattern<T>
): (data: T | undefined, error: Error | undefined, loading: boolean) => ReactNode {
  return (data: T | undefined, error: Error | undefined, loading: boolean): ReactNode => {
    return (
      <>
        {loading && p.Loading()}
        {error && p.Error(error)}
        {data && p.Data(data)}
      </>
    )
  }
}
