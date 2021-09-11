import React, { FC } from 'react'

import { ButtonIcon } from '../components'
import { useNavbar } from '../contexts/NavbarContext'
import styles from './Appbar.module.css'

interface Props {
  title?: string
}

export const Appbar: FC<Props> = ({ title, children }) => {
  const navbar = useNavbar()

  return (
    <div className={styles.appBar}>
      <ButtonIcon id="appbar-menu" icon="menu" onClick={() => (navbar.opened ? navbar.close() : navbar.open())} />
      <h1>{title}</h1>
      {children}
    </div>
  )
}
