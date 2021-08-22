import React, { FC } from 'react'

import { classNames } from '../helpers'
import classes from './Box.module.css'

interface Props {
  title: string
  className?: string
  variant?: 'default' | 'warning' | 'danger'
}

export const Box: FC<Props> = ({ children, className, title, variant = 'default' }) => (
  <fieldset className={classNames(classes.box, className, classes[variant])}>
    <legend>{title}</legend>
    {children}
  </fieldset>
)
