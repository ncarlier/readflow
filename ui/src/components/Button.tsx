import React, { ReactNode, ReactType } from 'react'
import Ink from 'react-ink'

import { classNames } from '../helpers'
import styles from './Button.module.css'
import Icon from './Icon'
import { PropsOf } from './PropsOf'

interface ButtonProps {
  icon?: string
  variant?: 'default' | 'primary' | 'danger' | 'flat'
  disabled?: boolean
  children: ReactNode
}

function Button<Tag extends ReactType = 'button'>(props: { as?: Tag } & ButtonProps & PropsOf<Tag>) {
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

export default Button
