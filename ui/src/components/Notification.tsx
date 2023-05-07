import React, { FC, PropsWithChildren } from 'react'

import { basicMarkdownToHTML, classNames } from '../helpers'
import styles from './Notification.module.css'

interface Props extends PropsWithChildren {
  message: string
  variant?: 'error' | 'info' | 'warning'
}

export const Notification: FC<Props> = ({ message, variant = 'info', children }) => {
  const className = classNames(styles.notification, styles[variant])

  return (
    <div className={className}>
      <div className={styles.message} dangerouslySetInnerHTML={{__html: basicMarkdownToHTML(message)}} />
      <div className={styles.actions}>{children}</div>
    </div>
  )
}
