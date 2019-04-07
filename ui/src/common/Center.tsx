import React, { ReactNode } from 'react'

import styles from './Center.module.css'
import { classNames } from './helpers'

type Props = {
  children: ReactNode
  className?: string
}

export default ({children, className}: Props) => (
  <section className={classNames(styles.center, className)}>
    <div>
      {children}
    </div>
  </section>
)
