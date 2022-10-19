import React, { FC, PropsWithChildren } from 'react'

import styles from './Header.module.css'

export const Header: FC<PropsWithChildren> = ({ children }) => <header className={styles.header}>{children}</header>
