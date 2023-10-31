import React, { FC, ImgHTMLAttributes } from 'react'

import { Article, ArticleThumbnail } from '../models'
import { getAPIURL } from '../../helpers'

const getThumbnailURL = (thumbnail: ArticleThumbnail, src: string) => `${getAPIURL()}/img/${thumbnail.hash}/resize:fit:${thumbnail.size}/${btoa(src)}`

const getThumbnailAttributes = ({thumbnails, image}: Article) => {
  const attrs :ImgHTMLAttributes<HTMLImageElement> = {}
  if (!thumbnails || thumbnails.length == 0) {
    return attrs
  }
  const sizes = thumbnails.reverse().map(thumb => `${thumb.size}px`)
  attrs.sizes = `(max-width: ${sizes[0]}) ${sizes.join(', ')}`
  attrs.srcSet = thumbnails?.map(thumb => `${getThumbnailURL(thumb, image)} ${thumb.size}w`).join(',')
  return attrs
}

interface Props {
  article: Article
}

export const ArticleImage: FC<Props> = ({ article }) => {
  let attrs :ImgHTMLAttributes<HTMLImageElement> = {}
  if (article.image && article.image.match(/^https?:\/\//)) {
    attrs = getThumbnailAttributes(article)
  }
  return (
    <img
      {...attrs}
      src={article.image}
      onError={(e) => (e.currentTarget.style.display = 'none')}
      crossOrigin='anonymous'
    />
  )
}
