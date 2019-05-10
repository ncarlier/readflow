import React, { ReactNode } from 'react'

import styles from './ArticleFooter.module.css'

interface Props {
  children?: ReactNode
}

type AllProps = Props

export default ({ children }: AllProps) => (
  <footer>
    <div className={styles.actions}>{children}</div>
  </footer>
)
