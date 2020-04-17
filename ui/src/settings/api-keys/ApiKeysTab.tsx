import React, { useContext, useState } from 'react'
import { useMutation, useQuery } from 'react-apollo-hooks'
import { useModal } from 'react-modal-hook'
import { RouteComponentProps } from 'react-router'
import { Link } from 'react-router-dom'

import Button from '../../components/Button'
import ConfirmDialog from '../../components/ConfirmDialog'
import DataTable, { OnSelectedFn } from '../../components/DataTable'
import Loader from '../../components/Loader'
import Masked from '../../components/Masked'
import Panel from '../../components/Panel'
import TimeAgo from '../../components/TimeAgo'
import { MessageContext } from '../../context/MessageContext'
import ErrorPanel from '../../error/ErrorPanel'
import { getGQLError, matchResponse } from '../../helpers'
import { usePageTitle } from '../../hooks'
import Bookmarklet from './Bookmarklet'
import { updateCacheAfterDelete } from './cache'
import { ApiKey, DeleteApiKeyRequest, DeleteApiKeyResponse, GetApiKeysResponse } from './models'
import { DeleteApiKeys, GetApiKeys } from './queries'

const definition = [
  {
    title: 'Title',
    render: (val: ApiKey) => (
      <Link title="Edit API key" to={`/settings/api-keys/${val.id}`}>
        {val.alias}
      </Link>
    )
  },
  {
    title: 'Token',
    render: (val: ApiKey) => <Masked value={val.token} />
  },
  {
    title: 'Bookmarklet',
    render: (val: ApiKey) => <Bookmarklet token={val.token} />
  },
  {
    title: 'Last usage',
    render: (val: ApiKey) => <TimeAgo dateTime={val.last_usage_at} />
  },
  {
    title: 'Created',
    render: (val: ApiKey) => <TimeAgo dateTime={val.created_at} />
  },
  {
    title: 'Updated',
    render: (val: ApiKey) => <TimeAgo dateTime={val.updated_at} />
  }
]

type AllProps = RouteComponentProps<{}>

export default ({ match }: AllProps) => {
  usePageTitle('Settings - API keys')

  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [selection, setSelection] = useState<number[]>([])
  const { data, error, loading } = useQuery<GetApiKeysResponse>(GetApiKeys)
  const [deleteApiKeysMutation] = useMutation<DeleteApiKeyResponse, DeleteApiKeyRequest>(DeleteApiKeys)
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
      if (res.data) {
        const nb = res.data.deleteAPIKeys
        showMessage(nb > 1 ? `${nb} API keys removed` : 'API key removed')
      }
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
    Data: data => <DataTable definition={definition} data={data.apiKeys} onSelected={onSelectedHandler} />,
    Other: () => <ErrorPanel>Unable to fetch API keys with no obvious reason :(</ErrorPanel>
  })

  return (
    <Panel>
      <header>
        {selection.length > 0 && (
          <Button id="remove-selection" title="Remove selection" variant="danger" onClick={showDeleteConfirmModal}>
            Remove
          </Button>
        )}
        <Button
          id="add-new-api-key"
          title="Add new API key"
          variant="primary"
          as={Link}
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
