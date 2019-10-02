import { LocationDescriptor } from 'history'
import React, { MouseEventHandler } from 'react'
import Ink from 'react-ink'
import { Link } from 'react-router-dom'

import styles from './ButtonIcon.module.css'
import { classNames } from '../helpers'
import Icon from './Icon'

interface Props {
  id?: string
  icon: string
  primary?: boolean
  loading?: boolean
  floating?: boolean
  title?: string
  to?: LocationDescriptor
  //onClick?: MouseEventHandler<HTMLAnchorElement|HTMLButtonElement>
  onClick?: MouseEventHandler
}

export default (props: Props) => {
  const { icon, primary, loading, floating, to, onClick, ...attrs } = props
  let className = classNames(styles.button, primary ? styles.primary : null, floating ? styles.floating : null)

  if (loading) {
    className = classNames(className, styles.loading)
    return (
      <button {...attrs} disabled className={className}>
        <Icon name="loop" />
      </button>
    )
  }

  if (to) {
    return (
      <Link {...attrs} to={to} className={className} onClick={onClick}>
        {icon && <Icon name={icon} />}
        <Ink />
      </Link>
    )
  }

  return (
    <button {...attrs} className={className} onClick={onClick}>
      <Icon name={icon} />
      <Ink />
    </button>
  )
}
