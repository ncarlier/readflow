import React from 'react'
import { format } from 'timeago.js'

type Props = {
  dateTime?: string
}

export default ({ dateTime }: Props) => {
  if (!dateTime) {
    return <span>-</span>
  }
  const date = new Date(dateTime)
  const ago = format(dateTime) 
  return (
    <time dateTime={date.toISOString()} title={date.toISOString()}>{ago}</time>
  )
}
