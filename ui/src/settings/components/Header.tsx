import React, {ReactNode} from 'react'

import styles from './Header.module.css'

type Props = {
  children: ReactNode
}

export default ({children}: Props) => (
  <header className={styles.header}>
    {children}
  </header>
)
