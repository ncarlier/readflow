import { ReactNode } from "react"
import { ApolloError } from "apollo-boost"
import { FormState } from "react-use-form-state"

export interface GQLResponsePattern<T> {
  Loading: () => ReactNode
  Error: (err: ApolloError | Error) => ReactNode
  Data: (data: T) => ReactNode
  Other: () => ReactNode
}

export function matchResponse<T>(p: GQLResponsePattern<T>): (data: T | undefined, error: ApolloError | Error | undefined, loading: boolean) => ReactNode {
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

export function classNames(...names: (string|undefined|null)[]) {
  return names.filter(name => name).join(' ')
}

export function getGQLError(err: any) {
  // console.log('Error', JSON.stringify(err, null, 4))
  if (err.networkError && err.networkError.name === 'ServerError') {
    return err.networkError.result.errors[0].message
  }
  return err.message
}

export function getURLParam<T>(params: URLSearchParams, name: string, fallback: T): T {
  let result = fallback
  if (params.has(name)) {
    const val = params.get(name)
    switch (typeof fallback) {
      case 'number':
        if (parseInt(val!, 10) != NaN) {
          result = <any>parseInt(val!, 10)
        }
        break
      case 'string':
        result = <any>val!
    }
  }
  return result
}

export function isValidForm(form: FormState<any>) {
  let result = true
  Object.keys(form.values).forEach(key => {
    if (!form.validity[key]) {
      result = false
    }
  })
  return result
}
