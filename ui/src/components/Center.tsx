import React, { FC, PropsWithChildren } from 'react'

import styles from './Center.module.css'
import { classNames } from '../helpers'

interface Props extends PropsWithChildren {
  className?: string
}

export const Center: FC<Props> = ({ children, className }) => (
  <section className={classNames(styles.center, className)}>
    <div>{children}</div>
  </section>
)
