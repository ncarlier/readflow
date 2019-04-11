import React, { ReactNode, useCallback } from 'react'
import Ink from 'react-ink'
import { Link } from 'react-router-dom'
import { connect } from 'react-redux'

import { RouterState } from 'connected-react-router'
import { ApplicationState } from '../store'

import styles from './LinkIcon.module.css'
import Icon from './Icon'
import { classNames } from './helpers'
import { LocationDescriptor } from 'history'

type IProps = {
  to?: LocationDescriptor
  icon: string
  active?: boolean
  onClick?: Function
  title?: string
  children: ReactNode
}

type Props = IProps & IPropsFromState

export const LinkIcon = (props: Props) => {
  const {children, icon, to, active, title, onClick} = props
  const { pathname } = props.router.location
  const handleOnClick = useCallback(
    () => onClick ? onClick() : () => true,
    [onClick]
  )

  let className = styles.link
  if (active || (to && pathname.startsWith(typeof to === 'string' ? to : to.pathname!))) {
    className = classNames(className, 'active')
  }

  if (!to) {
    return (
      <a style={{position: 'relative'}}
        onClick={handleOnClick}
        title={title}
        className={className}>
        <Icon name={icon}/>
        {children}
        <Ink />
      </a>
    )
  }

  return (
    <Link
      to={to}
      style={{position: 'relative'}}
      onClick={handleOnClick}
      title={title}
      className={className}
      >
      <Icon name={icon}/>
      {children}
      <Ink />
    </Link>
  )
}

interface IPropsFromState {
  router: RouterState
}

const mapStateToProps = ({ router }: ApplicationState): IPropsFromState => ({
  router
})

export default connect(
  mapStateToProps
)(LinkIcon)
