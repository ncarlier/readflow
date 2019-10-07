import React, { ReactNode } from 'react'

import { classNames } from '../helpers'
import classes from './Box.module.css'

interface Props {
  title: string
  children: ReactNode
  className?: string
  warning?: boolean
  danger?: boolean
}

export default ({ children, className, title, warning, danger }: Props) => (
  <fieldset className={classNames(classes.box, className, danger ? classes.danger : warning ? classes.warning : null)}>
    <legend>{title}</legend>
    {children}
  </fieldset>
)
