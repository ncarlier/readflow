
import React, { useCallback } from 'react'

import {Article} from '../models'

import LinkIcon from '../../common/LinkIcon'

type Props = {
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
    <LinkIcon
      title="Share article"
      onClick={share}
      icon="share">
      <span>Share</span>
    </LinkIcon>
  )
}
