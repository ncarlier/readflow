import React from 'react'

import useKeyboard, { KeyHandler } from '../hooks/useKeyboard'

interface Props {
  keys: string
  onKeypress: KeyHandler
}

export default ({ keys, onKeypress }: Props) => {
  useKeyboard(keys, onKeypress)

  return <kbd>{keys}</kbd>
}
