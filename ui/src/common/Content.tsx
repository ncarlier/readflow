import React, { ReactNode } from 'react'

import styles from './Content.module.css'

interface Props {
  id?: string
  children: ReactNode
}

export default ({ children, id }: Props) => (
  <section id={id} className={styles.content}>
    {children}
  </section>
)
