import React, { forwardRef, ReactNode, Ref } from 'react'

import { classNames } from '../helpers'
import styles from './Panel.module.css'

interface Props {
  children: ReactNode
  style?: any
  tabIndex?: number
  className?: string
}

export const Panel = forwardRef(({ children, className, ...rest }: Props, ref: Ref<any>) => (
  <section ref={ref} className={classNames(styles.panel, className)} {...rest}>
    {children}
  </section>
))
