import { ApolloError } from 'apollo-boost'
import React, { MouseEvent, ReactNode } from 'react'
import { FormState } from 'react-use-form-state'

import { API_BASE_URL } from '../constants'
import { FormMountValidity } from '../hooks/useOnMountInputValidator'

export const URLRegExp = new RegExp(
  '^(https?:\\/\\/)?' + // protocol
  '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.?)+[a-z]{2,}|' + // domain name
  '((\\d{1,3}\\.){3}\\d{1,3}))' + // OR ip (v4) address
  '(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*' + // port and path
  '(\\?[;&a-z\\d%_.~+=-]*)?' + // query string
    '(\\#[-a-z\\d_]*)?$',
  'i'
) // fragment locator

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

export function classNames(...names: (string | undefined | null)[]) {
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
    if (val && typeof fallback === 'number') {
      // eslint-disable-next-line use-isnan
      if (parseInt(val, 10) != NaN) {
        return Number.parseInt(val, 10) as any
      }
    }
    if (val && typeof fallback === 'string') {
      return val as any
    }
  }
  return result
}

export function isValidForm(formState: FormState<any>, onMountValidator: FormMountValidity<any>) {
  const validity = { ...onMountValidator.validity, ...formState.validity }
  Object.keys(formState.values).forEach(key => {
    if (!validity[key]) {
      return false
    }
  })
  return true
}

export function getBookmarklet(apiKey: string) {
  const { origin } = document.location
  const cred = btoa('api:' + apiKey)
  return `javascript:(function(){
FP_URL="${API_BASE_URL}";
FP_CRED="${cred}";
var js=document.body.appendChild(document.createElement("script"));
js.onerror=function(){alert("Sorry, unable to load bookmarklet.")};
js.src="${origin}/bookmarklet.js"})();`
}

export function preventBookmarkletClick(e: MouseEvent<any>) {
  e.preventDefault()
  alert("Don't click on me! But drag and drop me to your toolbar.")
}

export function getOnlineStatus() {
  return typeof navigator !== 'undefined' && typeof navigator.onLine === 'boolean' ? navigator.onLine : true
}
