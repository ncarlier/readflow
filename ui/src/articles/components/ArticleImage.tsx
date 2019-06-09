import React from 'react'

import { API_BASE_URL } from '../../constants'
import { useMedia } from '../../hooks'

interface Props {
  src: string
}

const proxifyImageURL = (url: string, width: number) => `${API_BASE_URL}/img?url=${encodeURIComponent(url)}&width=${width}`

export default ({ src }: Props) => {
  const mobileDisplay = useMedia('(max-width: 767px)')
  return (
    <img
      src={src}
      srcSet={`${proxifyImageURL(src, 320)} 320w,
              ${proxifyImageURL(src, 767)} 767w`}
      sizes={mobileDisplay ? '100vw' : '320px'}
      alt=""
      onError={e => (e.currentTarget.style.display = 'none')}
    />
  )
}
