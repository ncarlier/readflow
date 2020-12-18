import React, { useCallback } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { useModal } from 'react-modal-hook'

import Button from '../../../components/Button'
import ConfirmDialog from '../../../components/ConfirmDialog'
import { getGQLError } from '../../../helpers'
import { DeleteOutgoingWebhookRequest, DeleteOutgoingWebhookResponse } from './models'
import { DeleteOutgoingWebhooks } from './queries'

interface Props {
  selection: number[]
  onSuccess: (msg: string) => void
  onError: (err: string) => void
}

export default ({ selection, onSuccess, onError }: Props) => {
  const [deleteOutgoingWebhooksMutation] = useMutation<DeleteOutgoingWebhookResponse, DeleteOutgoingWebhookRequest>(
    DeleteOutgoingWebhooks
  )

  const deleteOutgoingWebhooks = useCallback(
    async (ids: number[], callback: () => void) => {
      try {
        const res = await deleteOutgoingWebhooksMutation({
          variables: { ids },
        })
        if (res.data) {
          const nb = res.data.deleteOutgoingWebhooks
          onSuccess(nb > 1 ? `${nb} outgoing webhooks removed` : 'Outgoing webhook removed')
        }
      } catch (err) {
        onError(getGQLError(err))
      } finally {
        callback()
      }
    },
    [onError, onSuccess, deleteOutgoingWebhooksMutation]
  )

  const [showDeleteConfirmModal, hideDeleteConfirmModal] = useModal(
    () => (
      <ConfirmDialog
        title="Delete outgoing webhook?"
        confirmLabel="Delete"
        onConfirm={() => deleteOutgoingWebhooks(selection, hideDeleteConfirmModal)}
        onCancel={hideDeleteConfirmModal}
      >
        Deleting an outgoing webhook is irreversible. Please confirm!
      </ConfirmDialog>
    ),
    [selection]
  )

  if (selection.length > 0) {
    return (
      <Button id="remove-selection-2" title="Remove selection" variant="danger" onClick={showDeleteConfirmModal}>
        Remove
      </Button>
    )
  }
  return null
}
