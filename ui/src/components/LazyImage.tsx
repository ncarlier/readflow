import React, { useState, useRef, useEffect, ImgHTMLAttributes, FC } from 'react'

import { useMedia } from '../hooks'
import { classNames, getAPIURL } from '../helpers'

import styles from './LazyImage.module.css'

const proxifyImageURL = (url: string, width: number) => getAPIURL(`/img?url=${encodeURIComponent(url)}&width=${width}`)

export const LazyImage: FC<ImgHTMLAttributes<HTMLImageElement>> = ({ src, ...attrs }) => {
  const mobileDisplay = useMedia('(max-width: 767px)')
  const [loaded, setLoaded] = useState(false)
  const imgRef = useRef<HTMLImageElement>(null)
  useEffect(() => {
    if (imgRef.current && imgRef.current.complete) {
      setLoaded(true)
    }
  }, [])

  if (!src) {
    return <img {...attrs} />
  }

  return (
    <div className={styles.wrapper}>
      <img {...attrs} src={proxifyImageURL(src, 64)} aria-hidden="true" className={styles.lqip} />
      <img
        ref={imgRef}
        {...attrs}
        src={src}
        srcSet={`${proxifyImageURL(src, 320)} 320w,
                ${proxifyImageURL(src, 767)} 767w`}
        sizes={mobileDisplay ? '100vw' : '320px'}
        loading="lazy"
        className={classNames(styles.source, loaded ? styles.loaded : null)}
        onLoad={() => setLoaded(true)}
        onError={(e) => (e.currentTarget.style.display = 'none')}
      />
    </div>
  )
}
