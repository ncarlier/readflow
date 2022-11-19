import React, { FC, PropsWithChildren } from 'react'

import styles from './DropDownMenuItem.module.css'

interface Props extends PropsWithChildren {
  label?: string
}

export const DropDownMenuItem: FC<Props> = (props) => {
  const { label, children } = props
  return (
    <div className={styles.menu_item}>
      <span>{label}</span>
      {children}
    </div>
  )
}
