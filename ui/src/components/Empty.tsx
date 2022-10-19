import React, { FC, PropsWithChildren } from 'react'

import { Center } from '.'
import styles from './Empty.module.css'

export const Empty: FC<PropsWithChildren> = ({ children }) => (
  <Center>
    <span className={styles.empty}>{children}</span>
  </Center>
)
