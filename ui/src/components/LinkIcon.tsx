import React, { ElementType, ReactNode } from 'react'
import Ink from 'react-ink'

import { classNames } from '../helpers'
import Icon from './Icon'
import styles from './LinkIcon.module.css'
import { PropsOf } from './PropsOf'

interface LinkIconProps {
  icon: string | ReactNode
  active?: boolean
  badge?: string | number
  children?: ReactNode
}

function LinkIcon<Tag extends ElementType = 'a'>(props: { as?: Tag } & LinkIconProps & PropsOf<Tag>) {
  const { as: Element = 'a', children, badge, icon, active, ...attrs } = props

  let className = styles.link
  if (active) {
    className = classNames(className, 'active')
  }

  return (
    <Element {...attrs} style={{ position: 'relative' }} className={className}>
      {typeof icon === 'string' ? <Icon name={icon} /> : icon}
      {children}
      {!!badge && <span className={styles.badge}>{badge}</span>}
      <Ink />
    </Element>
  )
}

export default LinkIcon
