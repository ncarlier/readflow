/* eslint-disable @typescript-eslint/no-non-null-assertion */
import React, { FormEvent, useRef, useState } from 'react'
import { Link } from 'react-router-dom'

import Empty from '../../components/Empty'
import TimeAgo from '../../components/TimeAgo'
import { ArchiveService } from './models'

export interface OnSelectedFn {
  (ids: number[]): void
}

interface Props {
  data: ArchiveService[]
  onSelected?: OnSelectedFn
}

export default ({ data, onSelected }: Props) => {
  const selectAllRef = useRef<HTMLInputElement>(null)
  const [selection, setSelection] = useState<Map<number, boolean>>(() => {
    const state = new Map<number, boolean>()
    data.forEach(service => service.id && state.set(service.id, false))
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
    return <Empty>No archive service yet</Empty>
  }

  return (
    <table>
      <thead>
        <tr>
          <th>
            <input ref={selectAllRef} type="checkbox" onChange={onCheckboxAllChange} />
          </th>
          <th>Alias</th>
          <th>Provider</th>
          <th>Created</th>
          <th>Updated</th>
        </tr>
      </thead>
      <tbody>
        {data.map(service => (
          <tr key={`archive-service-${service.id}`}>
            <th>
              <input type="checkbox" onChange={onCheckboxChange(service.id!)} checked={selection.get(service.id!)} />
            </th>
            <th>
              <Link title="Edit API key" to={`/settings/archive-services/${service.id}`}>
                {service.alias} {service.is_default && '(default)'}
              </Link>
            </th>
            <td>{service.provider}</td>
            <td>
              <TimeAgo dateTime={service.created_at} />
            </td>
            <td>
              <TimeAgo dateTime={service.updated_at} />
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
