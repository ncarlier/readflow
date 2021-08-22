import React, { FC } from 'react'

import styles from './Center.module.css'
import { classNames } from '../helpers'

interface Props {
  className?: string
}

export const Center: FC<Props> = ({ children, className }) => (
  <section className={classNames(styles.center, className)}>
    <div>{children}</div>
  </section>
)
