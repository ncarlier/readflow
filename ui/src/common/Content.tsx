import React, { ReactNode } from 'react'

import styles from './Content.module.css'

type Props = {
  children: ReactNode
}

export default ({children}: Props) => (
  <section className={styles.content}>
    {children}
  </section>
)
