import React, { CSSProperties, FunctionComponent } from 'react'

import { classNames } from '../helpers'
import styles from './Overlay.module.css'

interface Props {
  style?: CSSProperties
  className?: string
  visible: boolean
}

const Overlay: FunctionComponent<Props> = ({ children, className, visible, ...rest }) => {
  if (!visible) {
    return null
  }
  return (
    <section className={classNames(styles.overlay, className)} {...rest}>
      {children}
    </section>
  )
}

export default Overlay
