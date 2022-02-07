import gql from 'graphql-tag'
import React, { useContext } from 'react'
import { useMutation } from '@apollo/client'
import { useModal } from 'react-modal-hook'

import { Box, Button, ConfirmDialog } from '../../components'
import { useMessage } from '../../contexts'
import { getGQLError } from '../../helpers'
import { AuthContext } from 'react-oidc-context'

const DeleteAccount = gql`
  mutation {
    deleteAccount
  }
`

interface DeleteAccountResponse {
  deleteAccount: boolean
}

const DeleteAccountBox = () => {
  const auth = useContext(AuthContext)
  const { showErrorMessage } = useMessage()
  const [deleteAccountMutation] = useMutation<DeleteAccountResponse>(DeleteAccount)
  const deleteAccount = async () => {
    try {
      const res = await deleteAccountMutation()
      if (res.data && res.data.deleteAccount) {
        auth?.signoutRedirect()
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
    <Box title="Account" variant="danger">
      <p>When you delete your account, your profile, articles, settings will be permanently removed.</p>
      <Button title="Delete my account" onClick={showDeleteAccountModal}>
        Delete account
      </Button>
    </Box>
  )
}

export default DeleteAccountBox
