import React from 'react'
import { Kbd, ToggleMenuItem } from '../../components'
import { ArticleStatus } from '../models'

type ViewMode = ArticleStatus | 'starred'

interface Props {
  value: ViewMode
  onChange: (value: ViewMode) => void
  kbs: string
}

const values: {
  value: ViewMode,
  icon: string,
  title: string
}[] = [
  {
    value: 'inbox',
    icon: 'inbox',
    title: 'Inbox'
  },
  {
    value: 'to_read',
    icon: 'book',
    title: 'To read'
  },
  {
    value: 'read',
    icon: 'history',
    title: 'History'
  },
  {
    value: 'starred',
    icon: 'star',
    title: 'Starred'
  }
]

const toggle = (value: ViewMode) => value === 'inbox' ? 'to_read'
  : value === 'to_read' ? 'read'
  : value === 'read' ? 'starred'
  : 'inbox'

export const ToggleView = ({value, onChange, kbs}: Props) => (
  <>
    <ToggleMenuItem name='view-mode' value={value} onChange={onChange} values={values} />
    <Kbd keys={kbs} onKeypress={() => onChange(toggle(value))} />
  </>
)
