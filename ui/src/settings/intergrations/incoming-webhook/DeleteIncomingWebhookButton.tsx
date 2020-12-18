import React, { useCallback } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { useModal } from 'react-modal-hook'

import Button from '../../../components/Button'
import ConfirmDialog from '../../../components/ConfirmDialog'
import { getGQLError } from '../../../helpers'
import { updateCacheAfterDelete } from './cache'
import { DeleteIncomingWebhookRequest, DeleteIncomingWebhookResponse } from './models'
import { DeleteIncomingWebhooks } from './queries'

interface Props {
  selection: number[]
  onSuccess: (msg: string) => void
  onError: (err: string) => void
}

export default ({ selection, onSuccess, onError }: Props) => {
  const [deleteIncomingWebhooksMutation] = useMutation<DeleteIncomingWebhookResponse, DeleteIncomingWebhookRequest>(
    DeleteIncomingWebhooks
  )

  const deleteIncomingWebhooks = useCallback(
    async (ids: number[], callback: () => void) => {
      try {
        const res = await deleteIncomingWebhooksMutation({
          variables: { ids },
          update: updateCacheAfterDelete(ids),
        })
        if (res.data) {
          const nb = res.data.deleteIncomingWebhooks
          onSuccess(nb > 1 ? `${nb} incoming webhooks removed` : 'Incoming webhook removed')
        }
      } catch (err) {
        onError(getGQLError(err))
      } finally {
        callback()
      }
    },
    [onError, onSuccess, deleteIncomingWebhooksMutation]
  )

  const [showDeleteConfirmModal, hideDeleteConfirmModal] = useModal(
    () => (
      <ConfirmDialog
        title="Delete incoming webhook?"
        confirmLabel="Delete"
        onConfirm={() => deleteIncomingWebhooks(selection, hideDeleteConfirmModal)}
        onCancel={hideDeleteConfirmModal}
      >
        Deleting an incoming webhook is irreversible. Please confirm!
      </ConfirmDialog>
    ),
    [selection]
  )

  if (selection.length > 0) {
    return (
      <Button id="remove-selection-1" title="Remove selection" variant="danger" onClick={showDeleteConfirmModal}>
        Remove
      </Button>
    )
  }
  return null
}
