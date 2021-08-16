import React, { ReactNode, useContext } from 'react'
import CSSTransition from 'react-transition-group/CSSTransition'
import { useModal } from 'react-modal-hook'

import InfoDialog from '../components/InfoDialog'
import Shortcuts from '../components/Shortcuts'
import Snackbar from '../components/Snackbar'
import { NavbarContext } from '../contexts/NavbarContext'
import { classNames, isMobileDevice } from '../helpers'
import useKeyboard from '../hooks/useKeyboard'
import classes from './AppLayout.module.css'
import Navbar from './Navbar'
import UpdateAvailableNotification from '../components/UpdateAvailableNotification'
import useTheme from '../hooks/useTheme'

interface Props {
  children: ReactNode
}

export default (props: Props) => {
  const { children } = props

  // Activate theme
  useTheme()

  // Shortcuts global modal
  const [showShortcutsModal, hideShortcutsModal] = useModal(() => (
    <InfoDialog title="Shortcuts" onOk={hideShortcutsModal}>
      <Shortcuts />
    </InfoDialog>
  ))
  useKeyboard('?', showShortcutsModal)

  // const small = useMedia('(max-width: 400px)')
  // const large = useMedia('(min-width: 767px)')
  const navbar = useContext(NavbarContext)

  const deviceClassName = isMobileDevice() ? classes.mobile : null

  return (
    <div className={classNames(classes.layout, deviceClassName)}>
      <CSSTransition in={navbar.opened} classNames="fold" timeout={300} unmountOnExit>
        <aside>
          <Navbar />
        </aside>
      </CSSTransition>
      <section>
        {navbar.opened && <div id="navbar-fog" className={classes.fog} onClick={() => navbar.close()} />}
        {children}
        <UpdateAvailableNotification />
        <Snackbar />
      </section>
    </div>
  )
}
