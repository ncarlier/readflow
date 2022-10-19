import React, { ElementType, PropsWithChildren } from 'react'
import Ink from 'react-ink'

import { classNames } from '../helpers'
import styles from './Button.module.css'
import { Icon } from '.'
import { PropsOf } from './PropsOf'

interface ButtonProps extends PropsWithChildren {
  icon?: string
  variant?: 'default' | 'primary' | 'danger' | 'flat'
  disabled?: boolean
}

export function Button<Tag extends ElementType = 'button'>(props: { as?: Tag } & ButtonProps & PropsOf<Tag>) {
  const { as: Element = 'button', variant = 'default', disabled, icon, children, ...attrs } = props
  const className = classNames(styles.button, styles[variant])

  return (
    <Element className={className} {...attrs} data-test={`btn-${variant}`}>
      {icon && <Icon name={icon} />}
      {children}
      {!disabled && <Ink />}
    </Element>
  )
}
