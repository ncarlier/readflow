import React, {useCallback, MouseEventHandler} from 'react'
import Ink from 'react-ink'

import Icon from './Icon'

import styles from './ButtonIcon.module.css'
import { Link } from 'react-router-dom'
import { classNames } from './helpers'
import { LocationDescriptor } from 'history'

type Props = { 
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
  const { icon, title, primary, loading, floating, to, onClick } = props
  let className = classNames(
    styles.button,
    primary ? styles.primary : null,
    floating ? styles.floating : null,
  )

  if (loading) {
    className = classNames(className, styles.loading)
    return (
      <button
        title={title}
        disabled
        className={className}>
        <Icon name="loop" />
      </button>
    )
  }

  if (to) {
    return (
      <Link
        to={to}
        title={title}
        className={className}
        onClick={onClick}>
      {icon && <Icon name={icon} />}
      <Ink />
      </Link>
    )
  } 

  return (
    <button
      title={title} 
      className={className}
      onClick={onClick}>
      <Icon name={icon} /><Ink />
    </button>
  )
}
