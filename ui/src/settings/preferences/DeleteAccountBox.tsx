import gql from 'graphql-tag'
import React, { useContext } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { useModal } from 'react-modal-hook'

import auth from '../../auth'
import Box from '../../components/Box'
import Button from '../../components/Button'
import ConfirmDialog from '../../components/ConfirmDialog'
import { MessageContext } from '../../context/MessageContext'
import { getGQLError } from '../../helpers'

const DeleteAccount = gql`
  mutation {
    deleteAccount
  }
`

interface DeleteAccountResponse {
  deleteAccount: boolean
}

export default () => {
  const { showErrorMessage } = useContext(MessageContext)
  const deleteAccountMutation = useMutation<DeleteAccountResponse>(DeleteAccount)
  const deleteAccount = async () => {
    try {
      const res = await deleteAccountMutation()
      if (res.data && res.data.deleteAccount) {
        auth.logout()
      }
    } catch (err) {
      showErrorMessage(getGQLError(err))
    }
  }

  const [showDeleteAccountModal, hideDeleteAccountModal] = useModal(
    () => (
      <ConfirmDialog
        title="Delete your account?"
        confirmLabel="Goodbye!"
        onConfirm={() => deleteAccount()}
        onCancel={hideDeleteAccountModal}
      >
        DELETING THIS ACCOUNT IS IRREVERSIBLE. PLEASE CONFIRM!
      </ConfirmDialog>
    ),
    []
  )

  return (
    <Box title="Account" danger>
      <p>DANGER ZONE</p>
      <p>When you delete your account, your profile, articles, settings will be permanently removed.</p>
      <Button title="Delete my account" onClick={showDeleteAccountModal}>
        Delete account
      </Button>
    </Box>
  )
}
