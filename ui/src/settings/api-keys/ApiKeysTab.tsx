import React, { useContext, useState } from 'react'
import { useMutation, useQuery } from 'react-apollo-hooks'
import { useModal } from 'react-modal-hook'
import { RouteComponentProps } from 'react-router'

import Button from '../../common/Button'
import ConfirmDialog from '../../common/ConfirmDialog'
import { getGQLError, matchResponse } from '../../common/helpers'
import Loader from '../../common/Loader'
import Panel from '../../common/Panel'
import { MessageContext } from '../../context/MessageContext'
import ErrorPanel from '../../error/ErrorPanel'
import { usePageTitle } from '../../hooks'
import ApiKeysTable, { OnSelectedFn } from './ApiKeysTable'
import { updateCacheAfterDelete } from './cache'
import { GetApiKeysResponse } from './models'
import { DeleteApiKeys, GetApiKeys } from './queries'

type AllProps = RouteComponentProps<{}>

export default ({ match }: AllProps) => {
  usePageTitle('Settings - API keys')

  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [selection, setSelection] = useState<number[]>([])
  const { data, error, loading } = useQuery<GetApiKeysResponse>(GetApiKeys)
  const deleteApiKeysMutation = useMutation<{ ids: number[] }>(DeleteApiKeys)
  const { showMessage } = useContext(MessageContext)

  const onSelectedHandler: OnSelectedFn = keys => {
    setSelection(keys)
  }

  const deleteApiKeys = async (ids: number[]) => {
    try {
      const res = await deleteApiKeysMutation({
        variables: { ids },
        update: updateCacheAfterDelete(ids)
      })
      setSelection([])
      // console.log('API keys removed', res)
      const nb = res.data.deleteAPIKeys
      showMessage(nb > 1 ? `${nb} API keys removed` : 'API key removed')
    } catch (err) {
      setErrorMessage(getGQLError(err))
    }
    // eslint-disable-next-line @typescript-eslint/no-use-before-define
    hideDeleteConfirmModal()
  }

  const [showDeleteConfirmModal, hideDeleteConfirmModal] = useModal(
    () => (
      <ConfirmDialog
        title="Delete API keys?"
        confirmLabel="Delete"
        onConfirm={() => deleteApiKeys(selection)}
        onCancel={hideDeleteConfirmModal}
      >
        Deleting an API key is irreversible. Please confirm!
      </ConfirmDialog>
    ),
    [selection]
  )

  const render = matchResponse<GetApiKeysResponse>({
    Loading: () => <Loader />,
    Error: err => <ErrorPanel title="Unable to fetch API keys">{err.message}</ErrorPanel>,
    Data: data => <ApiKeysTable data={data.apiKeys} onSelected={onSelectedHandler} />,
    Other: () => <ErrorPanel>Unable to fetch API keys with no obvious reason :(</ErrorPanel>
  })

  return (
    <Panel>
      <header>
        {selection.length > 0 && (
          <Button id="remove-selection" title="Remove selection" danger onClick={showDeleteConfirmModal}>
            Remove
          </Button>
        )}
        <Button
          id="add-new-api-key"
          title="Add new API key"
          primary
          to={{
            pathname: match.path + '/add',
            state: { modal: true }
          }}
        >
          Add API key
        </Button>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to delete API key(s)">{errorMessage}</ErrorPanel>}
        {render(data, error, loading)}
      </section>
    </Panel>
  )
}
