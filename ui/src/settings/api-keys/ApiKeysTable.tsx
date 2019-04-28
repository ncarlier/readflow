import React, { FormEvent, useRef, useState } from 'react'
import { Link } from 'react-router-dom'

import Empty from '../../common/Empty'
import { getBookmarklet, preventBookmarkletClick } from '../../common/helpers'
import Icon from '../../common/Icon'
import TimeAgo from '../../common/TimeAgo'
import { ApiKey } from './models'

export interface OnSelectedFn {
  (ids: number[]): void
}

interface Props {
  data: ApiKey[]
  onSelected?: OnSelectedFn
}

export default ({ data, onSelected }: Props) => {
  const selectAllRef = useRef<HTMLInputElement>(null)
  const [selection, setSelection] = useState<Map<number, boolean>>(() => {
    const state = new Map<number, boolean>()
    data.forEach(apiKey => state.set(apiKey.id, false))
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
    return <Empty>No API key yet</Empty>
  }

  return (
    <table>
      <thead>
        <tr>
          <th>
            <input ref={selectAllRef} type="checkbox" onChange={onCheckboxAllChange} />
          </th>
          <th>Alias</th>
          <th>Token</th>
          <th>Bookmarklet</th>
          <th>Last usage</th>
          <th>Created</th>
          <th>Updated</th>
        </tr>
      </thead>
      <tbody>
        {data.map(apiKey => (
          <tr key={`api-key-${apiKey.id}`}>
            <th>
              <input type="checkbox" onChange={onCheckboxChange(apiKey.id)} checked={selection.get(apiKey.id)} />
            </th>
            <th>
              <Link title="Edit API key" to={`/settings/api-keys/${apiKey.id}`}>
                {apiKey.alias}
              </Link>
            </th>
            <td>
              <span className="masked">{apiKey.token}</span>
            </td>
            <td>
              <a title="Bookmark me!" href={getBookmarklet(apiKey.token)} onClick={preventBookmarkletClick}>
                <Icon name="bookmark" />
              </a>
            </td>
            <td>
              <TimeAgo dateTime={apiKey.last_usage_at} />
            </td>
            <td>
              <TimeAgo dateTime={apiKey.created_at} />
            </td>
            <td>
              <TimeAgo dateTime={apiKey.updated_at} />
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
