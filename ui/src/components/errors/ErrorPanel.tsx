import React, { FC, PropsWithChildren, ReactNode } from 'react'

import styles from './ErrorPanel.module.css'

interface Props extends PropsWithChildren {
  title?: string
  actions?: ReactNode
}

export const ErrorPanel: FC<Props> = ({ title = 'Oh snap!', children, actions }) => (
  <section className={styles.error}>
    <header>
      <h1>{title}</h1>
    </header>
    <div>{children}</div>
    {actions && <footer>{actions}</footer>}
  </section>
)
