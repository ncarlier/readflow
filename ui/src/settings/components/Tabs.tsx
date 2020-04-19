import React from 'react'
import { RouteComponentProps, withRouter } from 'react-router'
import { Link } from 'react-router-dom'

import LinkIcon from '../../components/LinkIcon'
import styles from './Tabs.module.css'

interface TabItem {
  key: string
  label: string
  icon: string
}

interface Props {
  items: TabItem[]
}

export default withRouter(({ location: { pathname }, items }: Props & RouteComponentProps) => (
  <nav className={styles.tabs}>
    <ul>
      {items.map((item) => (
        <li key={item.key}>
          <LinkIcon
            as={Link}
            to={`/settings/${item.key}`}
            title={item.label}
            icon={item.icon}
            active={pathname.startsWith(`/settings/${item.key}`)}
          >
            <span>{item.label}</span>
          </LinkIcon>
        </li>
      ))}
    </ul>
  </nav>
))
