import React, { CSSProperties } from 'react'

import { classNames } from '../helpers'
import { Spinner } from '.'
import styles from './Loader.module.css'

interface Props {
  blur?: boolean
  center?: boolean
  style?: CSSProperties
}

export const Loader = ({ blur, style = {}, center }: Props) => {
  if (center) {
    style.textAlign = 'center'
  }
  return (
    <section className={blur ? classNames(styles.overlay, styles.blur) : styles.overlay} style={style}>
      <Spinner />
    </section>
  )
}
