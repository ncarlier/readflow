import React, { CSSProperties } from 'react'

import { classNames } from '../helpers'
import styles from './Loader.module.css'
import Spinner from './Spinner'

interface Props {
  blur?: boolean
  center?: boolean
  style?: CSSProperties
}

export default ({ blur, style = {}, center }: Props) => {
  if (center) {
    style.textAlign = 'center'
  }
  return (
    <section className={blur ? classNames(styles.overlay, styles.blur) : styles.overlay} style={style}>
      <Spinner />
    </section>
  )
}
