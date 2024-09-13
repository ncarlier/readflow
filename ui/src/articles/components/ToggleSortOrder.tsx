import React from 'react'
import { Kbd, ToggleMenuItem } from '../../components'
import { SortOrder } from '../models'

interface Props {
  value: SortOrder
  onChange: (order: SortOrder) => void
  kbs: string
}

const values: {
  value: SortOrder,
  icon: string,
  title: string
}[] = [
  {
    value: 'asc',
    icon: 'arrow_upward',
    title: 'Ascending order'
  },
  {
    value: 'desc',
    icon: 'arrow_downward',
    title: 'Descending order'
  }
]

const toggle = (value: SortOrder) => value === 'asc' ? 'desc' : 'asc'

export const ToggleSortOrder = ({value, onChange, kbs}: Props) => (
  <>
    <ToggleMenuItem name='sort-order' value={value} onChange={onChange} values={values} />
    <Kbd keys={kbs} onKeypress={() => onChange(toggle(value))} />
  </>
)
