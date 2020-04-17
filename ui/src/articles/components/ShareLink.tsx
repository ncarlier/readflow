import React, { useCallback } from 'react'

import LinkIcon from '../../components/LinkIcon'
import { Article } from '../models'

interface Props {
  article: Article
  keyboard?: boolean
}

export default ({ article }: Props) => {
  const share = useCallback(() => {
    const nvg: any = window.navigator
    nvg.share({
      title: article.title,
      text: article.text,
      url: article.url
    })
  }, [article])

  return (
    <LinkIcon title="Share article" onClick={share} icon="share">
      <span>Share</span>
    </LinkIcon>
  )
}
