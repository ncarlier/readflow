import { RouterState } from 'connected-react-router'
import { LocationDescriptor } from 'history'
import React, { MouseEvent, ReactNode } from 'react'
import Ink from 'react-ink'
import { connect } from 'react-redux'
import { Link } from 'react-router-dom'

import { ApplicationState } from '../store'
import { classNames } from '../helpers'
import Icon from './Icon'
import styles from './LinkIcon.module.css'

interface IProps {
  id?: string
  to?: LocationDescriptor
  icon: string
  active?: boolean
  onClick?: (event: MouseEvent) => void
  title?: string
  badge?: string | number
  children: ReactNode
}

type Props = IProps & IPropsFromState

export const LinkIcon = (props: Props) => {
  const { children, badge, icon, to, active, router, ...attrs } = props
  const { pathname } = router.location
  const { id, title, onClick } = attrs

  let className = styles.link
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  if (active || (to && pathname.startsWith(typeof to === 'string' ? to : to.pathname!))) {
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

interface IPropsFromState {
  router: RouterState
}

const mapStateToProps = ({ router }: ApplicationState): IPropsFromState => ({
  router
})

export default connect(mapStateToProps)(LinkIcon)
