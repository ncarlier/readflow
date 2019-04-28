import React, { ReactNode } from 'react'

import styles from './Panel.module.css'

interface Props {
  children: ReactNode
  style?: any
  tabIndex?: number
}

export default ({ children, style, tabIndex }: Props) => (
  <section className={styles.panel} style={style} tabIndex={tabIndex}>
    {children}
  </section>
)
