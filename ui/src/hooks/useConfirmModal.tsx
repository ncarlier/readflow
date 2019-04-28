import React, { ReactNode } from 'react'
import { useModal } from 'react-modal-hook'

import ConfirmDialog from '../common/ConfirmDialog'

export default (title: string, body: ReactNode) => {
  const [showConfirmModal, hideConfirmModal] = useModal(() => (
    <ConfirmDialog title={title} confirmLabel="Ok" onConfirm={hideConfirmModal}>
      {body}
    </ConfirmDialog>
  ))
  return [showConfirmModal, hideConfirmModal]
}
