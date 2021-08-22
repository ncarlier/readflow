import React, { CSSProperties, FC } from 'react'

import { classNames } from '../helpers'
import styles from './Overlay.module.css'

interface Props {
  style?: CSSProperties
  className?: string
  visible: boolean
}

export const Overlay: FC<Props> = ({ children, className, visible, ...rest }) => {
  if (!visible) {
    return null
  }
  return (
    <section className={classNames(styles.overlay, className)} {...rest}>
      {children}
    </section>
  )
}
