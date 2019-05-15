import React, { ReactNode, useCallback, useState } from 'react'
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'
import { useModal } from 'react-modal-hook'

import { usePageTitle } from '../hooks'
import useKeyboard from '../hooks/useKeyboard'
import Appbar from './Appbar'
import Content from './Content'
import { isMobileDevice } from './device'
import { classNames } from './helpers'
import InfoDialog from './InfoDialog'
import Navbar from './Navbar'
import styles from './Page.module.css'
import Shortcuts from './Shortcuts'
import Snackbar from './Snackbar'

interface Props {
  title?: string
  subtitle?: string
  className?: string
  children: ReactNode
  actions?: ReactNode
}

export default (props: Props) => {
  const { children, title, subtitle, className, actions } = props

  usePageTitle(title, subtitle)

  const [showShortcutsModal, hideShortcutsModal] = useModal(() => (
    <InfoDialog title="Shortcuts" onOk={hideShortcutsModal}>
      <Shortcuts />
    </InfoDialog>
  ))
  useKeyboard('?', showShortcutsModal)

  // const small = useMedia('(max-width: 400px)')
  // const large = useMedia('(min-width: 767px)')
  const [navbarIsOpen, setNavbarIsOpen] = useState<boolean>(window.innerWidth > 767)
  const toggleNavbar = useCallback(() => setNavbarIsOpen(!navbarIsOpen), [navbarIsOpen])

  const deviceClassName = isMobileDevice() ? styles.mobile : null

  return (
    <div className={classNames(styles.page, className, deviceClassName)}>
      <ReactCSSTransitionGroup transitionName="fold" transitionEnterTimeout={300} transitionLeaveTimeout={300}>
        {navbarIsOpen && (
          <aside>
            <Navbar />
          </aside>
        )}
      </ReactCSSTransitionGroup>
      <section>
        {navbarIsOpen && <div id="navbar-fog" className={styles.fog} onClick={toggleNavbar} />}
        <Appbar title={title} onClickMenu={toggleNavbar} actions={actions} />
        <Content>{children}</Content>
        <Snackbar />
      </section>
    </div>
  )
}
