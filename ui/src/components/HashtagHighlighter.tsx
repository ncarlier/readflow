import React, { MouseEventHandler, useCallback } from 'react'

import styles from './HashtagHighlighter.module.css'
import { useHistory, useLocation } from 'react-router-dom'

type Props = {
  text: string
}

type Chunk = {
  type: 'hashtag' | 'text'
  value: string
}

const chunkTextWithHashtag = (text: string) => {
  const chunks = text.split(' ').reduce<Chunk[]>((acc, value) => {
    const type = value.startsWith('#') ? 'hashtag' : 'text'
    const last = acc.at(-1)
    if (last && type == 'text' && last.type == type) {
      // concat text
      acc.pop()
      value = last.value + ' ' + value
    }
    acc.push({type, value})
    return acc
  }, [])
  return chunks
}

export function HashtagHighlighter({ text }: Props) {
  const loc = useLocation()
  const { push } = useHistory()

  const handleClick = useCallback((value: string) => {
    const eventHandler: MouseEventHandler = function(event) {
      event.stopPropagation()
      event.preventDefault()
      const params = new URLSearchParams(loc.search)
      params.set('query', value)
      push({ ...loc, search: params.toString() })
    }
    return eventHandler
  }, [loc, push])

  const $chunks = chunkTextWithHashtag(text).map((chunk, idx) => {
    if (chunk.type == 'hashtag') {
      return <mark onClick={handleClick(chunk.value)} className={styles.hashtag} key={idx} title={`View other articles tagged with ${chunk.value}`}>
        {chunk.value}
      </mark>
    }
    return <span key={idx}>{chunk.value}</span>
  })
  return (<>{$chunks}</>)
}
