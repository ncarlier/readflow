import React, { FC } from 'react'
import CSSTransition from 'react-transition-group/CSSTransition'
import { useModal } from 'react-modal-hook'

import { InfoDialog, Shortcuts, Snackbar, UpdateAvailableNotification } from '../components'
import { useNavbar } from '../contexts/NavbarContext'
import { classNames, isMobileDevice } from '../helpers'
import { useKeyboard, useTheme } from '../hooks'
import classes from './AppLayout.module.css'
import { Navbar } from '.'

export const AppLayout: FC = ({ children }) => {
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
  const navbar = useNavbar()

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
