import React, { useState } from 'react'

import { LinkIcon, Loader } from '../../../components'
import { EditArticleLink } from './EditArticleLink'
import { Article } from '../../models'
import { RefreshArticleLink } from './RefreshArticleLink'

interface Props {
  article: Article
  keyboard?: boolean
  showEditModal?: () => void
  onCancel: (e: any) => void
}

export default (props: Props) => {
  const isOnline = !props.article.isOffline
  const hasUrl = props.article.url && props.article.url.trim() !== ''
  const [loading, setLoading] = useState(false)
  
  if (loading) {
    return <Loader blur />
  }

  return (
    <ul>
      {props.showEditModal && (
        <li>
          <EditArticleLink {...props} showEditModal={props.showEditModal}/>
        </li>
      )}
      {isOnline && hasUrl && (
        <li>
          <RefreshArticleLink {...props} onLoading={setLoading} />
        </li>
      )}
      <li>
        <LinkIcon title="Cancel" icon="cancel" onClick={props.onCancel}>
          <span>Cancel</span>
        </LinkIcon>
      </li>
    </ul>
  )
}
