import React, { ReactNode, useCallback, useState } from 'react'
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'
import { useModal } from 'react-modal-hook'

import InfoDialog from '../components/InfoDialog'
import Shortcuts from '../components/Shortcuts'
import Snackbar from '../components/Snackbar'
import { classNames, isMobileDevice } from '../helpers'
import useKeyboard from '../hooks/useKeyboard'
import classes from './AppLayout.module.css'
import Navbar from './Navbar'

interface Props {
  children: ReactNode
}

export default (props: Props) => {
  const { children } = props

  // Shortcuts global modal
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

  const deviceClassName = isMobileDevice() ? classes.mobile : null

  return (
    <div className={classNames(classes.layout, deviceClassName)}>
      <ReactCSSTransitionGroup transitionName="fold" transitionEnterTimeout={300} transitionLeaveTimeout={300}>
        {navbarIsOpen && (
          <aside>
            <Navbar />
          </aside>
        )}
      </ReactCSSTransitionGroup>
      <section>
        {navbarIsOpen && <div id="navbar-fog" className={classes.fog} onClick={toggleNavbar} />}
        {children}
        <Snackbar />
      </section>
    </div>
  )
}
