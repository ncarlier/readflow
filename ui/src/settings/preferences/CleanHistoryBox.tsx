import gql from 'graphql-tag'
import React from 'react'
import { useMutation } from '@apollo/client'
import { useModal } from 'react-modal-hook'

import { Category } from '../../categories/models'
import { Box, Button, ConfirmDialog } from '../../components'
import { useMessage } from '../../contexts'
import { getGQLError } from '../../helpers'

const CleanHistory = gql`
  mutation {
    cleanHistory {
      _inbox
    }
  }
`

interface CleanHistoryResponse {
  cleanHistory: Category[]
}

const CleanHistoryBox = () => {
  const { showMessage, showErrorMessage } = useMessage()
  const [cleanHistoryMutation] = useMutation<CleanHistoryResponse>(CleanHistory)
  const cleanHistory = async () => {
    try {
      const res = await cleanHistoryMutation()
      // console.log('Categories removed', res)
      if (res.data && res.data.cleanHistory) {
        showMessage('History cleaned')
      }
    } catch (err) {
      showErrorMessage(getGQLError(err))
    }
    // eslint-disable-next-line @typescript-eslint/no-use-before-define
    hideCleanHistoryModal()
  }

  const [showCleanHistoryModal, hideCleanHistoryModal] = useModal(
    () => (
      <ConfirmDialog
        title="Clean history?"
        confirmLabel="Delete"
        onConfirm={() => cleanHistory()}
        onCancel={hideCleanHistoryModal}
      >
        Deleting all read articles is irreversible. Please confirm!
      </ConfirmDialog>
    ),
    []
  )

  return (
    <Box title="History" variant="warning">
      <p>
        Read articles are kept temporarily and deleted eventually.
        <br />
        You can anticipate the deletion by cleaning the history.
      </p>
      <Button title="Clean the history" onClick={showCleanHistoryModal}>
        Clean history
      </Button>
    </Box>
  )
}

export default CleanHistoryBox
