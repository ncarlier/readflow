import React, { ReactNode, useContext } from 'react'

import ButtonIcon from '../components/ButtonIcon'
import { NavbarContext } from '../context/NavbarContext'
import styles from './Appbar.module.css'

interface Props {
  title?: string
  actions?: ReactNode
}

export default ({ title, actions }: Props) => {
  const navbar = useContext(NavbarContext)

  return (
    <div className={styles.appBar}>
      <ButtonIcon id="appbar-menu" icon="menu" onClick={() => (navbar.opened ? navbar.close() : navbar.open())} />
      <h1>{title}</h1>
      {actions}
    </div>
  )
}
