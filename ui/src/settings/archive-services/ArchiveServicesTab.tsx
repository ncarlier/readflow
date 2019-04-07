import React, { useState } from 'react'
import { usePageTitle } from '../../hooks'
import Panel from '../../common/Panel'
import Button from '../../common/Button'
import { useQuery, useMutation } from 'react-apollo-hooks'
import { matchResponse, getGQLError } from '../../common/helpers'
import Loader from '../../common/Loader'
import {OnSelectedFn} from './ArchiveServicesTable'
import { RouteComponentProps } from 'react-router'
import ErrorPanel from '../../error/ErrorPanel'
import { updateCacheAfterDelete } from './cache'
import { useModal } from 'react-modal-hook'
import ConfirmDialog from '../../common/ConfirmDialog'
import { GetArchiveServicesResponse } from './models'
import { GetArchiveServices, DeleteArchiveServices } from './queries'
import ArchiveServicesTable from './ArchiveServicesTable'
import { connectMessageDispatch, IMessageDispatchProps } from '../../containers/MessageContainer'

type AllProps = RouteComponentProps<{}> & IMessageDispatchProps

export const ArchiveServicesTab = ({match, showMessage}: AllProps) => {
  usePageTitle('Settings - Archive services')

  const [errorMessage, setErrorMessage] = useState<string | null>(null) 
  const [selection, setSelection] = useState<number[]>([])
  const { data, error, loading } = useQuery<GetArchiveServicesResponse>(GetArchiveServices)
  const deleteArchiveServicesMutation = useMutation<{ids: number[]}>(DeleteArchiveServices)

  const onSelectedHandler: OnSelectedFn = (keys) => {
    setSelection(keys)
  }

  const deleteArchiveServices = async (ids: number[]) => {
    try{
      const res = await deleteArchiveServicesMutation({
        variables: {ids},
        update: updateCacheAfterDelete(ids)
      })
      setSelection([])
      const nb = res.data.deleteArchivers
      showMessage(nb > 1 ? `${nb} archive services removed` : 'Archive service removed')
    } catch (err) {
      setErrorMessage(getGQLError(err))
    }
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
    Error: (err) => <ErrorPanel title="Unable to fetch archive services">
        {err.message}
      </ErrorPanel>,
    Data: ({archivers}) => <ArchiveServicesTable
      data={archivers}
      onSelected={onSelectedHandler}
    />,
    Other: () => <ErrorPanel>
        Unable to fetch archive services with no obvious reason :(
      </ErrorPanel>
  })

  return (
    <Panel>
      <header>
        { selection.length > 0 &&
          <Button
            title="Remove selection"
            danger
            onClick={showDeleteConfirmModal}>
            Remove
          </Button>
        }
        <Button 
          title="Add new archive service"
          primary
          to={{
            pathname: match.path + '/add',
            state: { modal: true }
          }} >
          Add archive service
        </Button>
      </header>
      <section>
        {errorMessage != null &&
          <ErrorPanel title="Unable to delete archive service(s)">
            {errorMessage}
          </ErrorPanel>
        }
        {render(data, error, loading)}
      </section>
    </Panel>
  )
}

export default connectMessageDispatch(ArchiveServicesTab)
