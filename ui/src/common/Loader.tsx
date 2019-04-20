import React from 'react'

import styles from './Loader.module.css'
import { classNames } from './helpers'

type Props = {
  blur?: boolean
}

export default ({ blur }: Props) => (
  <section className={blur ? classNames(styles.overlay, styles.blur) : styles.overlay}>
    <div className={styles.spinner} />
  </section>
)
