import React, { ReactNode } from 'react'

import Center from './Center'
import styles from './Empty.module.css'

interface Props {
  children: ReactNode
}

export default ({ children }: Props) => (
  <Center>
    <span className={styles.empty}>{children}</span>
  </Center>
)
