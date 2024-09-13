import React from 'react'
import { Kbd, ToggleMenuItem } from '../../components'
import { DisplayMode } from '../../contexts'

interface Props {
  value: DisplayMode
  onChange: (value: DisplayMode) => void
  kbs: string
}

const values: {
  value: DisplayMode,
  icon: string,
  title: string
}[] = [
  {
    value: 'grid',
    icon: 'dashboard',
    title: 'Display as grid'
  },
  {
    value: 'list',
    icon: 'list',
    title: 'Display as list'
  }
]

const toggle = (value: DisplayMode) => value === 'grid' ? 'list' : 'grid'

export const ToggleDisplayMode = ({value, onChange, kbs}: Props) => (
  <>
    <ToggleMenuItem name='display-mode' value={value} onChange={onChange} values={values} />
    <Kbd keys={kbs} onKeypress={() => onChange(toggle(value))} />
  </>
)
