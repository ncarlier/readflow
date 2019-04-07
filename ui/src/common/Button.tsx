import React, { ReactNode, RefObject, forwardRef, Ref } from 'react'
import Ink from 'react-ink'
import { Link } from 'react-router-dom'
import { LocationDescriptor } from 'history'

import Icon from './Icon'
import { classNames } from './helpers'

import styles from './Button.module.css'

type Props = { 
  icon?: string
  primary?: boolean
  danger?: boolean
  disabled?: boolean
  title?: string
  onClick?: (e: any) => void
  to?: LocationDescriptor
  children: ReactNode
}

export default forwardRef<any, Props>((props, ref) => {
  const {
    icon,
    title,
    primary,
    danger,
    disabled,
    onClick,
    to,
    children,
  } = props
  const className = classNames(
    styles.button,
    primary && !disabled ? styles.primary : undefined,
    danger && !disabled ? styles.danger : undefined,
  )

  if (to) {
    return (
      <Link
        ref={ref}
        to={to}
        title={title}
        className={className}
        onClick={onClick}
        >
      {icon && <Icon name={icon} />}
      {children}
      {!disabled && <Ink />}
      </Link>
    )
  } 
  return (
    <button
      ref={ref}
      title={title} 
      className={className}
      disabled={disabled}
      onClick={onClick}>
      {icon && <Icon name={icon} />}
      {children}
      {!disabled && <Ink />}
    </button>
  )
})
