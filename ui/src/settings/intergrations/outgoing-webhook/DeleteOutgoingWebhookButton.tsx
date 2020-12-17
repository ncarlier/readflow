import React, { useCallback } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { useModal } from 'react-modal-hook'

import Button from '../../../components/Button'
import ConfirmDialog from '../../../components/ConfirmDialog'
import { getGQLError } from '../../../helpers'
import { DeleteOutboundServiceRequest, DeleteOutboundServiceResponse } from './models'
import { DeleteOutboundServices } from './queries'

interface Props {
  selection: number[]
  onSuccess: (msg: string) => void
  onError: (err: string) => void
}

export default ({ selection, onSuccess, onError }: Props) => {
  const [deleteOutboundServicesMutation] = useMutation<DeleteOutboundServiceResponse, DeleteOutboundServiceRequest>(
    DeleteOutboundServices
  )

  const deleteOutboundServices = useCallback(
    async (ids: number[], callback: () => void) => {
      try {
        const res = await deleteOutboundServicesMutation({
          variables: { ids },
        })
        if (res.data) {
          const nb = res.data.deleteOutboundServices
          onSuccess(nb > 1 ? `${nb} outbound services removed` : 'Outbound service removed')
        }
      } catch (err) {
        onError(getGQLError(err))
      } finally {
        callback()
      }
    },
    [onError, onSuccess, deleteOutboundServicesMutation]
  )

  const [showDeleteConfirmModal, hideDeleteConfirmModal] = useModal(
    () => (
      <ConfirmDialog
        title="Delete outbound service?"
        confirmLabel="Delete"
        onConfirm={() => deleteOutboundServices(selection, hideDeleteConfirmModal)}
        onCancel={hideDeleteConfirmModal}
      >
        Deleting an outbound service is irreversible. Please confirm!
      </ConfirmDialog>
    ),
    [selection]
  )

  if (selection.length > 0) {
    return (
      <Button title="Remove selection" variant="danger" onClick={showDeleteConfirmModal}>
        Remove
      </Button>
    )
  }
  return null
}
