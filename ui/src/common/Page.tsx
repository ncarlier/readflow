import React, { ReactNode, useState, useCallback } from 'react'
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'

import Navbar from './Navbar'
import Appbar from './Appbar'
import Content from './Content'

import styles from './Page.module.css'

import { usePageTitle } from '../hooks'
import { classNames } from './helpers'
import Snackbar from './Snackbar'
import { useModal } from 'react-modal-hook'
import useKeyboard from '../hooks/useKeyboard'
import InfoDialog from './InfoDialog'
import Shortcuts from './Shortcuts'

type Props = {
  title?: string
  subtitle?: string
  className?: string
  children: ReactNode
  actions?: ReactNode
}

export default (props: Props) => {
  const {
    children,
    title,
    subtitle,
    className,
    actions,
  } = props

  usePageTitle(title, subtitle)

  const [showShortcutsModal, hideShortcutsModal] = useModal(
    () => (
      <InfoDialog
        title="Shortcuts"
        onOk={hideShortcutsModal}
      >
        <Shortcuts />
      </InfoDialog>
    )
  )
  useKeyboard('?', showShortcutsModal)

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
        <Appbar title={title} onClickMenu={toggleNavbar} actions={actions} />
        <Content>{children}</Content>
        <Snackbar />
      </section>
    </div>
  )
}
