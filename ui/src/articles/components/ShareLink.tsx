import React, { useCallback } from 'react'

import LinkIcon from '../../common/LinkIcon'
import { Article } from '../models'

interface Props {
  article: Article
}

export default ({ article }: Props) => {
  const nvg: any = window.navigator

  const share = useCallback(() => {
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
