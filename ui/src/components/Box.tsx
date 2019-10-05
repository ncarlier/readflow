import React, { ReactNode } from 'react'

import { classNames } from '../helpers'
import classes from './Box.module.css'

interface Props {
  title: string
  children: ReactNode
  className?: string
}

export default ({ children, className, title }: Props) => (
  <fieldset className={classNames(classes.box, className)}>
    <legend>{title}</legend>
    {children}
  </fieldset>
)
