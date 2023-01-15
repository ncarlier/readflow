import React, { SyntheticEvent, useCallback } from 'react'

import { LinkIcon } from '../../../components/LinkIcon'
import { Kbd } from '../../../components'

interface Props {
  keyboard?: boolean
  showEditModal: () => void
}

export const EditArticleLink = ({ keyboard, showEditModal }: Props) => {
  const handleOnClick = useCallback(
    (ev: SyntheticEvent) => {
      showEditModal()
      ev.preventDefault()
    },
    [showEditModal]
  )

  return (
    <LinkIcon title="Edit article" icon="edit" onClick={handleOnClick}>
      <span>Edit article ...</span>
      {keyboard && <Kbd keys="e" onKeypress={showEditModal} />}
    </LinkIcon>
  )
}
