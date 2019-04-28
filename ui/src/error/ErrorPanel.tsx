import React, { ReactNode } from 'react'

import styles from './ErrorPanel.module.css'

interface Props {
  title?: string
  children: ReactNode
  actions?: ReactNode
}

export default ({ title = 'Oh snap!', children, actions }: Props) => (
  <section className={styles.error}>
    <header>
      <h1>{title}</h1>
    </header>
    <div>{children}</div>
    {actions && <footer>{actions}</footer>}
  </section>
)
