/* eslint-disable @typescript-eslint/no-non-null-assertion */
import React from 'react'
import { RouteComponentProps, Link } from 'react-router-dom'

import ButtonIcon from '../components/ButtonIcon'
import { URLRegExp } from '../helpers'
import Page from '../layout/Page'
import AddArticleForm from './components/AddArticleForm'
import { Article } from './models'
import { Category } from '../categories/models'

interface Props {
  category?: Category
}

type AllProps = Props & RouteComponentProps

const Actions = () => <ButtonIcon as={Link} to="/unread" icon="arrow_back" title="back to the list" />

const extractURLFromParams = (qs: string) => {
  const params = new URLSearchParams(qs)
  if (params.has('url')) {
    return params.get('url')!
  } else if (params.has('text') || params.has('title')) {
    const text = params.get('text') || params.get('title')!
    if (URLRegExp.test(text)) {
      return text
    } else {
      const matches = text.match(/\bhttps?:\/\/\S+/gi)
      if (matches && URLRegExp.test(matches[0])) {
        return matches[0]
      }
    }
  }
  return undefined
}

export default ({ category, location, history }: AllProps) => {
  const url = extractURLFromParams(location.search)

  const onSuccess = (article: Article) => history.replace(`/unread/${article.id}`)
  const onCancel = () => history.replace('/unread')

  return (
    <Page title="Add new article" actions={<Actions />}>
      <AddArticleForm value={url} category={category} onCancel={onCancel} onSuccess={onSuccess} />
    </Page>
  )
}
