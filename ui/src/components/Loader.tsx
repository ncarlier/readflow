import React from 'react'

import { classNames } from '../helpers'
import styles from './Loader.module.css'

interface Props {
  blur?: boolean
}

export default ({ blur }: Props) => (
  <section className={blur ? classNames(styles.overlay, styles.blur) : styles.overlay}>
    <div className={styles.spinner} />
  </section>
)
