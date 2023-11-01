import React from 'react'
import { ago } from '../helpers'

interface Props {
  dateTime?: string
}

export const TimeAgo = ({ dateTime }: Props) => {
  if (!dateTime) {
    return <span>-</span>
  }
  const date = new Date(dateTime)
  return (
    <time dateTime={date.toISOString()} title={date.toISOString()}>
      {ago(date)}
    </time>
  )
}
