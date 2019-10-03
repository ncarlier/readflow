import { LocationDescriptor } from 'history'
import React, { MouseEvent, ReactNode } from 'react'
import Ink from 'react-ink'
import { Link } from 'react-router-dom'

import { classNames } from '../helpers'
import Icon from './Icon'
import styles from './LinkIcon.module.css'

interface Props {
  id?: string
  to?: LocationDescriptor
  icon: string
  active?: boolean
  onClick?: (event: MouseEvent) => void
  title?: string
  badge?: string | number
  children: ReactNode
}

export default (props: Props) => {
  const { children, badge, icon, to, active, ...attrs } = props
  const { id, title, onClick } = attrs

  let className = styles.link
  if (active) {
    className = classNames(className, 'active')
  }

  if (!to) {
    return (
      <a {...{ id, title, onClick }} style={{ position: 'relative' }} className={className}>
        <Icon name={icon} />
        {children}
        {!!badge && <span className={styles.badge}>{badge}</span>}
        <Ink />
      </a>
    )
  }

  return (
    <Link {...{ id, title, onClick }} to={to} style={{ position: 'relative' }} className={className}>
      <Icon name={icon} />
      {children}
      {!!badge && <span className={styles.badge}>{badge}</span>}
      <Ink />
    </Link>
  )
}
