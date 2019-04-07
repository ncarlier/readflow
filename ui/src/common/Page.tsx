import React, { ReactNode, useState, useCallback } from 'react'
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'

import Navbar from './Navbar'
import Appbar from './Appbar'
import Content from './Content'

import styles from './Page.module.css'

import { usePageTitle } from '../hooks'
import { classNames } from './helpers'
import Snackbar from './Snackbar'

type Props = {
  children: ReactNode
  title?: string
  subtitle?: string
  className?: string
  contextualMenu?: ReactNode
}

export default (props: Props) => {
  const {
    children,
    title = 'Reader',
    subtitle,
    className,
    contextualMenu,
  } = props

  usePageTitle(title, subtitle)

  // const small = useMedia('(max-width: 400px)')
  // const large = useMedia('(min-width: 767px)')
  const [navbarIsOpen, setNavbarIsOpen] = useState<boolean>(window.innerWidth > 767)
  const toggleNavbar = useCallback(() => setNavbarIsOpen(!navbarIsOpen), [navbarIsOpen])

  return (
    <div className={classNames(styles.page, className)}>
      <ReactCSSTransitionGroup
        transitionName="fold"
        transitionEnterTimeout={300}
        transitionLeaveTimeout={300}>
        {navbarIsOpen && <aside><Navbar /></aside>}
      </ReactCSSTransitionGroup>
      <section>
        { navbarIsOpen && <div className={styles.fog} onClick={toggleNavbar}/> }
        <Appbar title={title} onClickMenu={toggleNavbar} contextualMenu={contextualMenu} />
        <Content>{children}</Content>
        <Snackbar />
      </section>
    </div>
  )
}
