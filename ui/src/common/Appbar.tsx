import React, { useCallback, ReactNode } from 'react'

import ButtonIcon from './ButtonIcon'

import styles from './Appbar.module.css'
import DropdownMenu from './DropdownMenu';

type Props = {
  title?: string
  onClickMenu: Function
  contextualMenu?: ReactNode
}

export default ({title, onClickMenu, contextualMenu}: Props) => {
  const handleOnClickMenu = useCallback(
    () => onClickMenu(),
    [onClickMenu] 
  )

  return (
    <div className={styles.appBar}>
      <ButtonIcon icon="menu" onClick={handleOnClickMenu} />
      {title && <h1>{title}</h1>}
      {contextualMenu && <DropdownMenu>{contextualMenu}</DropdownMenu>}
    </div>
  )
}
