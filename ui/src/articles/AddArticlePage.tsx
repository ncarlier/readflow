/* eslint-disable @typescript-eslint/no-non-null-assertion */
import React from 'react'
import { RouteComponentProps, Link } from 'react-router-dom'

import ButtonIcon from '../components/ButtonIcon'
import { URLRegExp } from '../helpers'
import Page from '../layout/Page'
import AddArticleForm from './components/AddArticleForm'

type AllProps = RouteComponentProps

const Actions = () => <ButtonIcon as={Link} to="/unread" icon="arrow_back" title="back to the list" />

export default ({ location, history }: AllProps) => {
  const params = new URLSearchParams(location.search)

  let value = ''
  if (params.has('url')) {
    value = params.get('url')!
  } else if (params.has('text') || params.has('title')) {
    const text = params.get('text') || params.get('title')!
    if (URLRegExp.test(text)) {
      value = text
    } else {
      const matches = text.match(/\bhttps?:\/\/\S+/gi)
      if (matches && URLRegExp.test(matches[0])) {
        value = matches[0]
      }
    }
  }

  const redirect = () => history.replace('/unread')

  return (
    <Page title="Add new article" actions={<Actions />}>
      <AddArticleForm value={value} onCancel={redirect} onSuccess={redirect} />
    </Page>
  )
}
