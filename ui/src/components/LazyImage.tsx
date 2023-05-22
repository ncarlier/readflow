import React, { useState, useRef, useEffect } from 'react'

import { useMedia } from '../hooks'
import { classNames, getAPIURL } from '../helpers'

import styles from './LazyImage.module.css'

interface Props {
  src: string
  alt?: string
}

const proxifyImageURL = (url: string, width: number) => getAPIURL(`/img?url=${encodeURIComponent(url)}&width=${width}`)

export const LazyImage = ({ src, alt = '' }: Props) => {
  const mobileDisplay = useMedia('(max-width: 767px)')
  const [loaded, setLoaded] = useState(false)
  const imgRef = useRef<HTMLImageElement>(null)
  useEffect(() => {
    if (imgRef.current && imgRef.current.complete) {
      setLoaded(true)
    }
  }, [])

  return (
    <div className={styles.wrapper}>
      <img src={proxifyImageURL(src, 64)} aria-hidden="true" alt={alt} className={styles.lqip} />
      <img
        ref={imgRef}
        src={src}
        srcSet={`${proxifyImageURL(src, 320)} 320w,
                ${proxifyImageURL(src, 767)} 767w`}
        sizes={mobileDisplay ? '100vw' : '320px'}
        alt={alt}
        loading="lazy"
        className={classNames(styles.source, loaded ? styles.loaded : null)}
        onLoad={() => setLoaded(true)}
        onError={(e) => (e.currentTarget.style.display = 'none')}
      />
    </div>
  )
}
