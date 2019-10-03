import React, { ReactNode } from 'react'

import ButtonIcon from '../components/ButtonIcon'
import styles from './Appbar.module.css'

interface Props {
  title?: string
  actions?: ReactNode
}

export default ({ title, actions }: Props) => {
  const handleOnClickMenu = () => true

  return (
    <div className={styles.appBar}>
      <ButtonIcon id="appbar-menu" icon="menu" onClick={handleOnClickMenu} />
      <h1>{title}</h1>
      {actions}
    </div>
  )
}
