import React, { FormEvent, ReactNode, useRef, useState } from 'react'

import classes from './DataTable.module.css'
import Empty from './Empty'

export interface OnSelectedFn {
  (ids: number[]): void
}

interface DataTableDefinition {
  title: string
  render: (val: any) => ReactNode
}

interface RowValue {
  id?: number
  [prop: string]: any
}

interface Props {
  definition: DataTableDefinition[]
  data: RowValue[]
  onSelected?: OnSelectedFn
}

export default ({ definition, data, onSelected }: Props) => {
  const selectAllRef = useRef<HTMLInputElement>(null)
  const [selection, setSelection] = useState<Map<number, boolean>>(() => {
    const state = new Map<number, boolean>()
    data.forEach((val) => val.id && state.set(val.id, false))
    return state
  })

  const triggerOnSelectedEvent = (state: Map<number, boolean>) => {
    if (onSelected) {
      const payload = Array.from(state)
        .map((tuple) => {
          const [key, val] = tuple
          return val ? key : -1
        })
        .filter((v) => v !== -1)
      onSelected(payload)
    }
  }

  const onCheckboxChange = (id?: number) => (e: FormEvent<HTMLInputElement>) => {
    if (!id) {
      return
    }
    const newValue = e.currentTarget.checked
    const newState = new Map<number, boolean>(selection).set(id, newValue)
    const node = selectAllRef.current
    if (node) {
      let allChecked = true
      newState.forEach((val) => (val ? null : (allChecked = false)))
      node.checked = allChecked
    }
    triggerOnSelectedEvent(newState)
    setSelection(newState)
  }

  const onCheckboxAllChange = (e: FormEvent<HTMLInputElement>) => {
    const newValue = e.currentTarget.checked
    const newState = new Map<number, boolean>()
    selection.forEach((val, key) => newState.set(key, newValue))
    triggerOnSelectedEvent(newState)
    setSelection(newState)
  }

  if (data.length === 0) {
    return <Empty>No record</Empty>
  }

  return (
    <table className={classes.data_table}>
      <thead>
        <tr>
          <th>
            <input ref={selectAllRef} type="checkbox" onChange={onCheckboxAllChange} />
          </th>
          {definition.map((def) => (
            <th key={`dt-title-${def.title.toLowerCase()}`}>{def.title}</th>
          ))}
        </tr>
      </thead>
      <tbody>
        {data.map((val) => (
          <tr key={`dt-row-${val.id}`}>
            <th>
              <input type="checkbox" onChange={onCheckboxChange(val.id)} checked={!!val.id && selection.get(val.id)} />
            </th>
            {definition.map((def) => (
              <td key={`dt-row-${val.id}-${def.title.toLowerCase()}`}>{def.render(val)}</td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  )
}
