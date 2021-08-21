import React from 'react'

import { API_BASE_URL } from '../../constants'
import { useMedia } from '../../hooks'

interface Props {
  src: string
  alt?: string
}

const proxifyImageURL = (url: string, width: number) =>
  `${API_BASE_URL}/img?url=${encodeURIComponent(url)}&width=${width}`

export default ({ src, alt = '' }: Props) => {
  const mobileDisplay = useMedia('(max-width: 767px)')
  const attrs: any = {}
  if (src.match(/^https?:\/\//)) {
    attrs.srcSet = `${proxifyImageURL(src, 320)} 320w, ${proxifyImageURL(src, 767)} 767w`
  }
  return (
    <img
      src={src}
      {...attrs}
      sizes={mobileDisplay ? '100vw' : '320px'}
      alt={alt}
      onError={(e) => (e.currentTarget.style.display = 'none')}
    />
  )
}
