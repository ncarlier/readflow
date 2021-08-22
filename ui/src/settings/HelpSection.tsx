import React, { FC } from 'react'
import classes from './HelpSection.module.css'

const HelpSection: FC = ({ children }) => <div className={classes.help}>{children}</div>

export default HelpSection
