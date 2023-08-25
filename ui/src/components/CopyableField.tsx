import React, { useCallback } from 'react'
import { ButtonIcon } from '.'

interface Props {
  value: string
  masked?: boolean
}

export const CopyableField = ({ value, masked }: Props) => {
  const copyToClipboard = useCallback(() => navigator.clipboard.writeText(value), [value])

  return (
    <div>
      <span className={masked ? 'masked' : ''} >{value}</span>
      <ButtonIcon icon="file_copy" onClick={copyToClipboard} title="Copy to the clipboard" />
    </div>
  )
}
