import React, { ReactNode, useCallback } from 'react'

import styles from './Appbar.module.css'
import ButtonIcon from './ButtonIcon'

interface Props {
  title?: string
  onClickMenu: Function
  actions?: ReactNode
}

export default ({ title, onClickMenu, actions }: Props) => {
  const handleOnClickMenu = useCallback(() => onClickMenu(), [onClickMenu])

  return (
    <div className={styles.appBar}>
      <ButtonIcon icon="menu" onClick={handleOnClickMenu} />
      <h1>{title}</h1>
      {actions}
    </div>
  )
}
