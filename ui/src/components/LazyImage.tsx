import React, { useState, useRef, useEffect, ImgHTMLAttributes, FC } from 'react'

import { classNames, thumbHashToDataURL } from '../helpers'

import styles from './LazyImage.module.css'

const base64ToBinary = (b64: string) => new Uint8Array(window.atob(b64).split('').map(x => x.charCodeAt(0)))

interface Props {
  thumbhash: string
}

const hideFn = (ev: React.SyntheticEvent<HTMLElement, Event>) => {ev.currentTarget.style.display = 'none'}

export const LazyImage: FC<ImgHTMLAttributes<HTMLImageElement> & Props> = ({thumbhash, srcSet, ...attrs }) => {
  const [loaded, setLoaded] = useState(false)
  const [data, setData] = useState('')
  const [width, setWidth] = useState('0px')
  const [aspectRatio, setAspectRatio] = useState<number>()
  const imgRef = useRef<HTMLImageElement>(null)
  const lqipRef = useRef<HTMLImageElement>(null)
  useEffect(() => {
    if (imgRef.current && imgRef.current.complete) {
      setLoaded(true)
    }
  }, [])

  useEffect(() => {
    if (!thumbhash) {
      return
    }
    const [size, hash] = thumbhash.split('|')
    const [width, height] = size.split('x')
    setWidth(`${width}px`)
    try {
      if (height) {
        setAspectRatio(parseInt(width) / parseInt(height))
      }
      setData(thumbHashToDataURL(base64ToBinary(hash)))
    } catch (err) {
      console.error('unable to decode thumbhash', err)
    }
  }, [thumbhash])

  if (!thumbhash) {
    return <img {...attrs} onError={hideFn} />
  }

  return (
    <div className={styles.wrapper}>
      <img {...attrs}
        ref={lqipRef}
        width={width}
        src={data}
        aria-hidden="true"
        onAnimationEnd={hideFn}
        style={{ aspectRatio }}
        className={classNames(styles.lqip, loaded ? styles.loaded : null)}
      />
      <img
        ref={imgRef}
        {...attrs}
        srcSet={srcSet}
        loading="lazy"
        className={classNames(styles.source, loaded ? styles.loaded : null)}
        onLoad={() => setLoaded(true)}
        onError={hideFn}
      />
    </div>
  )
}
