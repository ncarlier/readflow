import React, { useCallback, FormEvent } from 'react'
import { Location } from 'history'
import { useLocation, useHistory } from 'react-router-dom'
import { useFormState } from 'react-use-form-state'

import Icon from '../../components/Icon'

import styles from './Search.module.css'
import { GetArticlesRequest } from '../models'

function getLocationWithQueryParam(loc: Location, query: string) {
  const params = new URLSearchParams(loc.search)
  params.set('query', query)
  return { ...loc, search: params.toString() }
}

interface Props {
  req: GetArticlesRequest
}

interface SearchFormFields {
  query: string
}

export default ({ req }: Props) => {
  const loc = useLocation()
  const { push } = useHistory()
  const [formState, { search }] = useFormState<SearchFormFields>({ query: req.query || '' })

  const handleOnSubmit = useCallback(
    (e: FormEvent) => {
      e.preventDefault()
      const { query } = formState.values
      push(getLocationWithQueryParam(loc, query))
      return false
    },
    [formState, loc, push]
  )

  return (
    <div className={styles.box}>
      <div className={styles.icon}>
        <Icon name="search" />
      </div>
      <form onSubmit={handleOnSubmit}>
        <input {...search('query')} placeholder="Search..." autoFocus={!!req.query} />
      </form>
    </div>
  )
}
