import React from 'react'
import { RouteComponentProps } from 'react-router-dom'

import Page from  '../common/Page'
import { URLRegExp } from '../common/helpers'
import ButtonIcon from '../common/ButtonIcon'
import AddArticleForm from './components/AddArticleForm'

type AllProps = RouteComponentProps

export default ({ location, history }: AllProps) => {
  const params = new URLSearchParams(location.search)

  let value = ""
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

  const redirect = () => history.replace("/unread")

  return (
    <Page title="Add new article"
          actions={
            <ButtonIcon
              to="/unread" 
              icon="arrow_back"
              title="back to the list"
            />
          }>
      <AddArticleForm value={value} onCancel={redirect} onSuccess={redirect} />
    </Page>
  )
}
