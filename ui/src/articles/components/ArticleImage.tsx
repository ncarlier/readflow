import React, { FC, ImgHTMLAttributes, useEffect, useState } from 'react'

import { Article, ArticleThumbnail } from '../models'
import { getAPIURL } from '../../helpers'
import { LazyImage } from '../../components'

const getThumbnailURL = (thumbnail: ArticleThumbnail, src: string) => `${getAPIURL()}/img/${thumbnail.hash}/resize:fit:${thumbnail.size}/${btoa(src)}`

const getThumbnailAttributes = (article: Article) => {
  const attrs :ImgHTMLAttributes<HTMLImageElement> = {}
  if (!article.thumbnails || article.thumbnails.length == 0) {
    return attrs
  }

  const thumbnails = [...article.thumbnails].sort((a, b) => parseInt(b.size) - parseInt(a.size))
  const sizes = thumbnails.map(thumb => `${thumb.size}px`)
  attrs.sizes = `(max-width: ${sizes[0]}) ${sizes.join(', ')}`
  attrs.srcSet = thumbnails.reverse().map(thumb => `${getThumbnailURL(thumb, article.image)} ${thumb.size}w`).join(',')
  return attrs
}

interface Props {
  article: Article
}

export const ArticleImage: FC<Props> = ({ article }) => {
  const [attrs, setAttrs] = useState<ImgHTMLAttributes<HTMLImageElement>>({})
  useEffect(() => {
    if (article.image && article.image.match(/^https?:\/\//)) {
      try {
        setAttrs(getThumbnailAttributes(article))
      } catch (err) {
        console.error('unable to get article thumbnail attributes', article, err)
      }
    }
  }, [article])
  
  return (
    <LazyImage
      {...attrs}
      thumbhash={article.thumbhash}
      src={article.image}
      // crossOrigin='anonymous'
    />
  )
}
