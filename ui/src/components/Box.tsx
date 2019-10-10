import React, { ReactNode } from 'react'

import { classNames } from '../helpers'
import classes from './Box.module.css'

interface Props {
  title: string
  children: ReactNode
  className?: string
  variant?: 'default' | 'warning' | 'danger'
}

export default ({ children, className, title, variant = 'default' }: Props) => (
  <fieldset className={classNames(classes.box, className, classes[variant])}>
    <legend>{title}</legend>
    {children}
  </fieldset>
)
