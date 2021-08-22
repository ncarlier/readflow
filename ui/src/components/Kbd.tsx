import React from 'react'

import { useKeyboard, KeyHandler } from '../hooks'

interface Props {
  keys: string
  onKeypress: KeyHandler
}

export const Kbd = ({ keys, onKeypress }: Props) => {
  useKeyboard(keys, onKeypress)

  return <kbd>{keys}</kbd>
}
