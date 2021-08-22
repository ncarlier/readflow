import React, { FC } from 'react'

import { classNames } from '../helpers'
import styles from './Notification.module.css'

interface Props {
  message: string
  variant?: 'error' | 'info' | 'warning'
}

export const Notification: FC<Props> = ({ message, variant = 'info', children }) => {
  const className = classNames(styles.notification, styles[variant])

  return (
    <div className={className}>
      <div className={styles.message}>{message}</div>
      <div className={styles.actions}>{children}</div>
    </div>
  )
}
