import React, { FormEvent, useCallback, useContext, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { RouteComponentProps } from 'react-router'
import { Link } from 'react-router-dom'
import { useFormState } from 'react-use-form-state'

import Button from '../../../components/Button'
import FormInputField from '../../../components/FormInputField'
import Panel from '../../../components/Panel'
import { MessageContext } from '../../../context/MessageContext'
import ErrorPanel from '../../../error/ErrorPanel'
import { getGQLError, isValidForm } from '../../../helpers'
import { usePageTitle } from '../../../hooks'
import { updateCacheAfterCreate } from './cache'
import IncomingWebhookHelp from './IncomingWebhookHelp'
import { CreateOrUpdateIncomingWebhookRequest, CreateOrUpdateIncomingWebhookResponse } from './models'
import { CreateOrUpdateIncomingWebhook } from './queries'

interface AddIncomingWebhookFormFields {
  alias: string
}

type AllProps = RouteComponentProps<{}>

export default ({ history }: AllProps) => {
  usePageTitle('Settings - Add new incoming webhook')

  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text }] = useFormState<AddIncomingWebhookFormFields>()
  const [addIncomingWebhookMutation] = useMutation<
    CreateOrUpdateIncomingWebhookResponse,
    CreateOrUpdateIncomingWebhookRequest
  >(CreateOrUpdateIncomingWebhook)
  const { showMessage } = useContext(MessageContext)

  const addIncomingWebhook = useCallback(
    async (incomingWebhook: CreateOrUpdateIncomingWebhookRequest) => {
      try {
        const res = await addIncomingWebhookMutation({
          variables: incomingWebhook,
          update: updateCacheAfterCreate,
        })
        if (res.data) {
          showMessage(`New incoming webhook: ${res.data.createOrUpdateIncomingWebhook.alias}`)
          // console.log('New incoming webhook', res)
          history.goBack()
        }
      } catch (err) {
        setErrorMessage(getGQLError(err))
      }
    },
    [addIncomingWebhookMutation, showMessage, history]
  )

  const handleOnSubmit = useCallback(
    (e: FormEvent | MouseEvent) => {
      e.preventDefault()
      if (!isValidForm(formState)) {
        setErrorMessage('Please fill out correctly the mandatory fields.')
        return
      }
      const { alias } = formState.values
      addIncomingWebhook({ alias })
    },
    [formState, addIncomingWebhook]
  )

  return (
    <Panel>
      <header>
        <h1>Add new incoming webhook</h1>
      </header>
      <section>
        <IncomingWebhookHelp />
        {errorMessage != null && <ErrorPanel title="Unable to add new incoming webhook">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField label="Alias" {...text('alias')} error={formState.errors.alias} required autoFocus />
        </form>
      </section>
      <footer>
        <Button title="Back to integrations" as={Link} to="/settings/integrations">
          Cancel
        </Button>
        <Button title="Add incoming webhook" onClick={handleOnSubmit} variant="primary">
          Add
        </Button>
      </footer>
    </Panel>
  )
}
