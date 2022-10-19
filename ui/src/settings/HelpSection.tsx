import React, { FC, PropsWithChildren } from 'react'
import classes from './HelpSection.module.css'

const HelpSection: FC<PropsWithChildren> = ({ children }) => <div className={classes.help}>{children}</div>

export default HelpSection
