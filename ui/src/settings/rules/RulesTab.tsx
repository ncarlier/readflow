import React, { useState } from 'react'
import { usePageTitle } from '../../hooks'
import Panel from '../../common/Panel'
import Button from '../../common/Button'
import { useQuery, useMutation } from 'react-apollo-hooks'
import { matchResponse, getGQLError } from '../../common/helpers'
import Loader from '../../common/Loader'
import { RouteComponentProps } from 'react-router'
import ErrorPanel from '../../error/ErrorPanel'
import { updateCacheAfterDelete } from './cache'
import { useModal } from 'react-modal-hook'
import ConfirmDialog from '../../common/ConfirmDialog'
import { GetRulesResponse } from './models'
import { GetRules, DeleteRules } from './queries'
import RulesTable, {OnSelectedFn} from './RulesTable'
import { connectMessageDispatch, IMessageDispatchProps } from '../../containers/MessageContainer';

type AllProps = RouteComponentProps<{}> & IMessageDispatchProps

export const RulesTab = ({match, showMessage}: AllProps) => {
  usePageTitle('Settings - Rules')

  const [errorMessage, setErrorMessage] = useState<string | null>(null) 
  const [selection, setSelection] = useState<number[]>([])
  const { data, error, loading } = useQuery<GetRulesResponse>(GetRules)
  const deleteRulesMutation = useMutation<{ids: number[]}>(DeleteRules)

  const onSelectedHandler: OnSelectedFn = (keys) => {
    setSelection(keys)
  }

  const deleteRules = async (ids: number[]) => {
    try{
      const res = await deleteRulesMutation({
        variables: {ids},
        update: updateCacheAfterDelete(ids)
      })
      setSelection([])
      const nb = res.data.deleteArchivers
      showMessage(nb > 1 ? `${nb} rules removed` : 'Rule removed')
    } catch (err) {
      setErrorMessage(getGQLError(err))
    }
    hideDeleteConfirmModal()
  }
  
  const [showDeleteConfirmModal, hideDeleteConfirmModal] = useModal(
    () => (
      <ConfirmDialog
        title="Delete rule?"
        confirmLabel="Delete"
        onConfirm={() => deleteRules(selection)}
        onCancel={hideDeleteConfirmModal}
      >
        Deleting a rule is irreversible. Please confirm!
      </ConfirmDialog>
    ),
    [selection]
  )

  const render = matchResponse<GetRulesResponse>({
    Loading: () => <Loader />,
    Error: (err) => <ErrorPanel title="Unable to fetch rules">
        {err.message}
      </ErrorPanel>,
    Data: ({rules}) => <RulesTable
      data={rules}
      onSelected={onSelectedHandler}
    />,
    Other: () => <ErrorPanel>
        Unable to fetch rules with no obvious reason :(
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
          title="Add new rule"
          primary
          to={{
            pathname: match.path + '/add',
            state: { modal: true }
          }} >
          Add rule
        </Button>
      </header>
      <section>
        {errorMessage != null &&
          <ErrorPanel title="Unable to delete rule(s)">
            {errorMessage}
          </ErrorPanel>
        }
        {render(data, error, loading)}
      </section>
    </Panel>
  )
}

export default connectMessageDispatch(RulesTab)
