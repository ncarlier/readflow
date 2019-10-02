import { LocationDescriptor } from 'history'
import React, { forwardRef, ReactNode } from 'react'
import Ink from 'react-ink'
import { Link } from 'react-router-dom'

import styles from './Button.module.css'
import { classNames } from '../helpers'
import Icon from './Icon'

interface Props {
  id?: string
  icon?: string
  primary?: boolean
  danger?: boolean
  disabled?: boolean
  autoFocus?: boolean
  title?: string
  onClick?: (e: any) => void
  to?: LocationDescriptor
  children: ReactNode
}

export default forwardRef<any, Props>((props, ref) => {
  // eslint-disable-next-line react/prop-types
  const { icon, primary, danger, disabled, to, children, ...attrs } = props
  const className = classNames(
    styles.button,
    primary && !disabled ? styles.primary : undefined,
    danger && !disabled ? styles.danger : undefined
  )
  const dataTest = primary ? 'btn-primary' : 'btn'

  if (to) {
    return (
      <Link ref={ref} to={to} className={className} {...attrs} data-test={dataTest}>
        {icon && <Icon name={icon} />}
        {children}
        {!disabled && <Ink />}
      </Link>
    )
  }
  return (
    <button ref={ref} className={className} disabled={disabled} {...attrs} data-test={dataTest}>
      {icon && <Icon name={icon} />}
      {children}
      {!disabled && <Ink />}
    </button>
  )
})
