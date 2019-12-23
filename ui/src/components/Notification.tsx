import React, { ReactNode } from 'react'

import { classNames } from '../helpers'
import styles from './Notification.module.css'

interface Props {
  children?: ReactNode
  message: string
  variant?: 'error' | 'info' | 'warning'
}

export default ({ message, variant = 'info', children }: Props) => {
  const className = classNames(styles.notification, styles[variant])

  return (
    <div className={className}>
      <div className={styles.message}>{message}</div>
      <div className={styles.actions}>{children}</div>
    </div>
  )
}
