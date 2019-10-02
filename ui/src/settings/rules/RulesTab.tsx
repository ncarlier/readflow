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
import { updateCacheAfterDelete } from './cache'
import { GetRulesResponse } from './models'
import { DeleteRules, GetRules } from './queries'
import RulesTable, { OnSelectedFn } from './RulesTable'

type AllProps = RouteComponentProps<{}>

export default ({ match }: AllProps) => {
  usePageTitle('Settings - Rules')

  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [selection, setSelection] = useState<number[]>([])
  const { data, error, loading } = useQuery<GetRulesResponse>(GetRules)
  const deleteRulesMutation = useMutation<{ ids: number[] }>(DeleteRules)
  const { showMessage } = useContext(MessageContext)

  const onSelectedHandler: OnSelectedFn = keys => {
    setSelection(keys)
  }

  const deleteRules = async (ids: number[]) => {
    try {
      const res = await deleteRulesMutation({
        variables: { ids },
        update: updateCacheAfterDelete(ids)
      })
      setSelection([])
      const nb = res.data.deleteArchivers
      showMessage(nb > 1 ? `${nb} rules removed` : 'Rule removed')
    } catch (err) {
      setErrorMessage(getGQLError(err))
    }
    // eslint-disable-next-line @typescript-eslint/no-use-before-define
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
    Error: err => <ErrorPanel title="Unable to fetch rules">{err.message}</ErrorPanel>,
    Data: data => <RulesTable data={data.rules} onSelected={onSelectedHandler} />,
    Other: () => <ErrorPanel>Unable to fetch rules with no obvious reason :(</ErrorPanel>
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
          title="Add new rule"
          primary
          to={{
            pathname: match.path + '/add',
            state: { modal: true }
          }}
        >
          Add rule
        </Button>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to delete rule(s)">{errorMessage}</ErrorPanel>}
        {render(data, error, loading)}
      </section>
    </Panel>
  )
}
