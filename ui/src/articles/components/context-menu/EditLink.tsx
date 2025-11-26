import React, { useState } from 'react'

import { LinkIcon, Overlay } from '../../../components'
import { Article } from '../../models'
import EditPanel from './EditPanel'

interface Props {
  article: Article
  keyboard?: boolean
  showEditModal?: () => void
}

export default (props: Props) => {
  const [isVisible, setIsVisible] = useState(false)
  const showOverlay = () => setIsVisible(true)
  const hideOverlay = () => setIsVisible(false)

  return (
    <>
      <LinkIcon title="Edit article ..." icon="edit" onClick={showOverlay}>
        <span>Edit ...</span>
      </LinkIcon>
      <Overlay visible={isVisible}>
        <EditPanel onCancel={hideOverlay} { ...props } />
      </Overlay>
    </>
  )
}
