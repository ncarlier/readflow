import React from 'react'
import { Icon } from './Icon'

import styles from './ToggleMenuItem.module.css'

interface ToggleMenuItemType {
  value: string
  icon: string
  title: string
}

interface Props {
  name: string
  value: string
  onChange: (value: any) => void
  values: ToggleMenuItemType[] 
}

export const ToggleMenuItem = ({ name, value, values, onChange }: Props) => (
  <ul className={styles.toggle_menu_item}>
    {values.map((item) => (
      <li key={`${name}-${item.value}`}>
        <input
          type='radio'
          id={`${name}-${item.value}`}
          name={name}
          checked={value === item.value}
          onChange={() => onChange(item.value)}
        />
        <label htmlFor={`${name}-${item.value}`} title={item.title}>
          <Icon name={item.icon} />
        </label>
      </li>
    ))}
  </ul>
)
