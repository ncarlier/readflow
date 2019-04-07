import React, { useState, useCallback } from 'react'
import { useModal } from "react-modal-hook"
import { usePageTitle } from '../../hooks'
import Panel from '../../common/Panel'
import Button from '../../common/Button'
import { useQuery, useMutation } from 'react-apollo-hooks'
import { matchResponse, getGQLError } from '../../common/helpers'
import Loader from '../../common/Loader'
import CategoriesTable, {OnSelectedFn} from './CategoriesTable'
import { RouteComponentProps } from 'react-router'
import { GetCategoriesResponse, DeleteCategoriesResponse } from '../../categories/models'
import { GetCategories, DeleteCategories } from '../../categories/queries'
import * as messageActions from '../../store/message/actions'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'
import ErrorPanel from '../../error/ErrorPanel'
import { updateCacheAfterDelete } from '../../categories/cache'
import ConfirmDialog from '../../common/ConfirmDialog';

type AllProps = RouteComponentProps<{}> & IPropsFromDispatch

export const CategoriesTab = ({match, showMessage}: AllProps) => {
  usePageTitle('Settings - Categories')

  const [errorMessage, setErrorMessage] = useState<string | null>(null) 
  const [selection, setSelection] = useState<number[]>([])
  const { data, error, loading } = useQuery<GetCategoriesResponse>(GetCategories)
  const deleteCategoriesMutation = useMutation<{ids: number[]}>(DeleteCategories)

  const onSelectedHandler: OnSelectedFn = (keys) => {
    setSelection(keys)
  }

  const deleteCategories = async (ids: number[]) => {
    try{
      const res = await deleteCategoriesMutation({
        variables: {ids},
        update: updateCacheAfterDelete(ids)
      })
      setSelection([])
      // console.log('Categories removed', res)
      const nb = res.data.deleteCategories
      showMessage(nb > 1 ? `${nb} categories removed` : 'Category removed')
    } catch (err) {
      setErrorMessage(getGQLError(err))
    }
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
    Error: (err) => <ErrorPanel title="Unable to fetch categories">
        {err.message}
      </ErrorPanel>,
    Data: ({categories}) => <CategoriesTable
      data={categories}
      onSelected={onSelectedHandler}
    />,
    Other: () => <ErrorPanel>
        Unable to fetch categories with no obvious reason :(
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
          title="Add new category"
          primary
          to={{
            pathname: match.path + '/add',
            state: { modal: true }
          }} >
          Add category
        </Button>
      </header>
      <section>
        {errorMessage != null &&
          <ErrorPanel title="Unable to delete categories">
            {errorMessage}
          </ErrorPanel>
        }
        {render(data, error, loading)}
      </section>
    </Panel>
  )
}

interface IPropsFromDispatch {
  showMessage: typeof messageActions.showMessage
}

const mapDispatchToProps = (dispatch: Dispatch): IPropsFromDispatch => ({
  showMessage: (msg: string|null) => dispatch(messageActions.showMessage(msg))
})

export default connect(
  null,
  mapDispatchToProps
)(CategoriesTab)
