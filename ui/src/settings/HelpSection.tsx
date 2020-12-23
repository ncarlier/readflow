import React, { ReactNode } from 'react'
import classes from './HelpSection.module.css'

interface Props {
  children?: ReactNode
}

export default ({ children }: Props) => <div className={classes.help}>{children}</div>
