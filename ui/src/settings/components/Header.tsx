import React, { FC } from 'react'

import styles from './Header.module.css'

export const Header: FC = ({ children }) => <header className={styles.header}>{children}</header>
