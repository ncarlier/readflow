import React, { FC } from 'react'

import { Center } from '.'
import styles from './Empty.module.css'

export const Empty: FC = ({ children }) => (
  <Center>
    <span className={styles.empty}>{children}</span>
  </Center>
)
