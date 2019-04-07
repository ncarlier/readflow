import React, { ReactNode } from 'react'

import styles from './HelpLink.module.css'
import Icon from './Icon'

type IProps = {
  href: string
  title?: string
  children: ReactNode
}

type Props = IProps

export default ({href, title, children}: Props) => (
  <a className={styles.help} href={href} title={title} target="_blank">
    <Icon name="help"/> {children}
  </a>
)
