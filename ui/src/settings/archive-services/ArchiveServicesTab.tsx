import React, { useContext, useState } from 'react'
import { useMutation, useQuery } from 'react-apollo-hooks'
import { useModal } from 'react-modal-hook'
import { RouteComponentProps } from 'react-router'

import Button from '../../components/Button'
import ConfirmDialog from '../../components/ConfirmDialog'
import { getGQLError, matchResponse } from '../../helpers'
import Loader from '../../components/Loader'
import Panel from '../../components/Panel'
import { MessageContext } from '../../context/MessageContext'
import ErrorPanel from '../../error/ErrorPanel'
import { usePageTitle } from '../../hooks'
import ArchiveServicesTable, { OnSelectedFn } from './ArchiveServicesTable'
import { updateCacheAfterDelete } from './cache'
import { GetArchiveServicesResponse } from './models'
import { DeleteArchiveServices, GetArchiveServices } from './queries'

type AllProps = RouteComponentProps<{}>

export default ({ match }: AllProps) => {
  usePageTitle('Settings - Archive services')

  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [selection, setSelection] = useState<number[]>([])
  const { data, error, loading } = useQuery<GetArchiveServicesResponse>(GetArchiveServices)
  const deleteArchiveServicesMutation = useMutation<{ ids: number[] }>(DeleteArchiveServices)
  const { showMessage } = useContext(MessageContext)

  const onSelectedHandler: OnSelectedFn = keys => {
    setSelection(keys)
  }

  const deleteArchiveServices = async (ids: number[]) => {
    try {
      const res = await deleteArchiveServicesMutation({
        variables: { ids },
        update: updateCacheAfterDelete(ids)
      })
      setSelection([])
      const nb = res.data.deleteArchivers
      showMessage(nb > 1 ? `${nb} archive services removed` : 'Archive service removed')
    } catch (err) {
      setErrorMessage(getGQLError(err))
    }
    // eslint-disable-next-line @typescript-eslint/no-use-before-define
    hideDeleteConfirmModal()
  }

  const [showDeleteConfirmModal, hideDeleteConfirmModal] = useModal(
    () => (
      <ConfirmDialog
        title="Delete archive service?"
        confirmLabel="Delete"
        onConfirm={() => deleteArchiveServices(selection)}
        onCancel={hideDeleteConfirmModal}
      >
        Deleting an archive service is irreversible. Please confirm!
      </ConfirmDialog>
    ),
    [selection]
  )

  const render = matchResponse<GetArchiveServicesResponse>({
    Loading: () => <Loader />,
    Error: err => <ErrorPanel title="Unable to fetch archive services">{err.message}</ErrorPanel>,
    Data: data => <ArchiveServicesTable data={data.archivers} onSelected={onSelectedHandler} />,
    Other: () => <ErrorPanel>Unable to fetch archive services with no obvious reason :(</ErrorPanel>
  })

  return (
    <Panel>
      <header>
        {selection.length > 0 && (
          <Button title="Remove selection" danger onClick={showDeleteConfirmModal}>
            Remove
          </Button>
        )}
        <Button
          title="Add new archive service"
          primary
          to={{
            pathname: match.path + '/add',
            state: { modal: true }
          }}
        >
          Add archive service
        </Button>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to delete archive service(s)">{errorMessage}</ErrorPanel>}
        {render(data, error, loading)}
      </section>
    </Panel>
  )
}
