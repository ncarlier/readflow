import React from 'react'

import styles from './Tabs.module.css'
import LinkIcon from '../../common/LinkIcon'

type TabItem = {
  key: string
  label: string
  icon: string
}

type Props = {
  items: TabItem[]
}

export default ({ items }: Props) => (
  <nav className={styles.tabs}>
    <ul>
      {items.map(item => (
        <li key={item.key}>
          <LinkIcon
            to={`/settings/${item.key}`}
            icon={item.icon}>
            {item.label}
          </LinkIcon>
        </li>
      ))}
    </ul>
  </nav>
)
