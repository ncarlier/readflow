import React, { useContext, useState } from 'react'
import { useMutation, useQuery } from 'react-apollo-hooks'
import { useModal } from 'react-modal-hook'
import { RouteComponentProps } from 'react-router'

import { updateCacheAfterDelete } from '../../categories/cache'
import { GetCategoriesResponse, DeleteCategoriesResponse, DeleteCategoriesRequest } from '../../categories/models'
import { DeleteCategories, GetCategories } from '../../categories/queries'
import Button from '../../components/Button'
import ConfirmDialog from '../../components/ConfirmDialog'
import { getGQLError, matchResponse } from '../..//helpers'
import Loader from '../../components/Loader'
import Panel from '../../components/Panel'
import { MessageContext } from '../../context/MessageContext'
import ErrorPanel from '../../error/ErrorPanel'
import { usePageTitle } from '../../hooks'
import CategoriesTable, { OnSelectedFn } from './CategoriesTable'

type AllProps = RouteComponentProps<{}>

export default ({ match }: AllProps) => {
  usePageTitle('Settings - Categories')

  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [selection, setSelection] = useState<number[]>([])
  const { data, error, loading } = useQuery<GetCategoriesResponse>(GetCategories)
  const deleteCategoriesMutation = useMutation<DeleteCategoriesResponse, DeleteCategoriesRequest>(DeleteCategories)
  const { showMessage } = useContext(MessageContext)

  const onSelectedHandler: OnSelectedFn = keys => {
    setSelection(keys)
  }

  const deleteCategories = async (ids: number[]) => {
    try {
      const res = await deleteCategoriesMutation({
        variables: { ids },
        update: updateCacheAfterDelete(ids)
      })
      setSelection([])
      // console.log('Categories removed', res)
      const nb = res.data!.deleteCategories
      showMessage(nb > 1 ? `${nb} categories removed` : 'Category removed')
    } catch (err) {
      setErrorMessage(getGQLError(err))
    }
    // eslint-disable-next-line @typescript-eslint/no-use-before-define
    hideDeleteConfirmModal()
  }

  const [showDeleteConfirmModal, hideDeleteConfirmModal] = useModal(
    () => (
      <ConfirmDialog
        title="Delete categories?"
        confirmLabel="Delete"
        onConfirm={() => deleteCategories(selection)}
        onCancel={hideDeleteConfirmModal}
      >
        Deleting a category is irreversible. Please confirm!
      </ConfirmDialog>
    ),
    [selection]
  )

  const render = matchResponse<GetCategoriesResponse>({
    Loading: () => <Loader />,
    Error: err => <ErrorPanel title="Unable to fetch categories">{err.message}</ErrorPanel>,
    Data: data => <CategoriesTable data={data.categories.filter(c => c.id !== null)} onSelected={onSelectedHandler} />,
    Other: () => <ErrorPanel>Unable to fetch categories with no obvious reason :(</ErrorPanel>
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
          id="add-new-category"
          title="Add new category"
          primary
          to={{
            pathname: match.path + '/add',
            state: { modal: true }
          }}
        >
          Add category
        </Button>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to delete categories">{errorMessage}</ErrorPanel>}
        {render(data, error, loading)}
      </section>
    </Panel>
  )
}
