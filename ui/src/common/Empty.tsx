import React, { ReactNode } from 'react'

import styles from './Empty.module.css'
import Center from './Center'

type Props = {
  children: ReactNode
}

export default ({children}: Props) => (
  <Center>
    <span className={styles.empty}>{children}</span>
  </Center>
)
