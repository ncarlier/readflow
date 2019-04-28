import React, { ReactNode } from 'react'

import styles from './Header.module.css'

interface Props {
  children: ReactNode
}

export default ({ children }: Props) => <header className={styles.header}>{children}</header>
