import React from 'react'

import LinkIcon from '../../common/LinkIcon'
import styles from './Tabs.module.css'

interface TabItem {
  key: string
  label: string
  icon: string
}

interface Props {
  items: TabItem[]
}

export default ({ items }: Props) => (
  <nav className={styles.tabs}>
    <ul>
      {items.map(item => (
        <li key={item.key}>
          <LinkIcon to={`/settings/${item.key}`} title={item.label} icon={item.icon}>
            <span>{item.label}</span>
          </LinkIcon>
        </li>
      ))}
    </ul>
  </nav>
)
