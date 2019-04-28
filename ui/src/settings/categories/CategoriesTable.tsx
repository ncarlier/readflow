/* eslint-disable @typescript-eslint/no-non-null-assertion */
import React, { FormEvent, useRef, useState } from 'react'
import { Link } from 'react-router-dom'

import { Category } from '../../categories/models'
import Empty from '../../common/Empty'
import TimeAgo from '../../common/TimeAgo'

export interface OnSelectedFn {
  (ids: number[]): void
}

interface Props {
  data: Category[]
  onSelected?: OnSelectedFn
}

export default ({ data, onSelected }: Props) => {
  const selectAllRef = useRef<HTMLInputElement>(null)
  const [selection, setSelection] = useState<Map<number, boolean>>(() => {
    const state = new Map<number, boolean>()
    data.forEach(cat => cat.id && state.set(cat.id, false))
    return state
  })

  const triggerOnSelectedEvent = (state: Map<number, boolean>) => {
    if (onSelected) {
      const payload = Array.from(state)
        .map(tuple => {
          const [key, val] = tuple
          return val ? key : -1
        })
        .filter(v => v !== -1)
      onSelected(payload)
    }
  }

  const onCheckboxChange = (id: number) => (e: FormEvent<HTMLInputElement>) => {
    const newValue = e.currentTarget.checked
    const newState = new Map<number, boolean>(selection).set(id, newValue)
    const node = selectAllRef.current
    if (node) {
      let allChecked = true
      newState.forEach(val => (val ? null : (allChecked = false)))
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
    return <Empty>No category yet</Empty>
  }

  return (
    <table>
      <thead>
        <tr>
          <th>
            <input ref={selectAllRef} type="checkbox" onChange={onCheckboxAllChange} />
          </th>
          <th>Title</th>
          <th>Created</th>
          <th>Updated</th>
        </tr>
      </thead>
      <tbody>
        {data.map(category => (
          <tr key={`cat-${category.id}`}>
            <th>
              <input type="checkbox" onChange={onCheckboxChange(category.id!)} checked={selection.get(category.id!)} />
            </th>
            <th>
              <Link title="Edit category" to={`/settings/categories/${category.id}`}>
                {category.title}
              </Link>
            </th>
            <td>
              <TimeAgo dateTime={category.created_at} />
            </td>
            <td>
              <TimeAgo dateTime={category.updated_at} />
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
