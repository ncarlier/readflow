import { History } from 'history'
import React, { FormEvent, useCallback, useContext, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { Link } from 'react-router-dom'
import { useFormState } from 'react-use-form-state'

import Button from '../../../components/Button'
import FormInputField from '../../../components/FormInputField'
import { MessageContext } from '../../../context/MessageContext'
import ErrorPanel from '../../../error/ErrorPanel'
import { getGQLError, isValidForm } from '../../../helpers'
import IncomingWebhookHelp from './IncomingWebhookHelp'
import { IncomingWebhook, CreateOrUpdateIncomingWebhookRequest, CreateOrUpdateIncomingWebhookResponse } from './models'
import { CreateOrUpdateIncomingWebhook } from './queries'

interface EditIncomingWebhookFormFields {
  alias: string
}

interface Props {
  data: IncomingWebhook
  history: History
}

export default ({ data, history }: Props) => {
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text }] = useFormState<EditIncomingWebhookFormFields>({ alias: data.alias })
  const [editIncomingWebhookMutation] = useMutation<
    CreateOrUpdateIncomingWebhookResponse,
    CreateOrUpdateIncomingWebhookRequest
  >(CreateOrUpdateIncomingWebhook)
  const { showMessage } = useContext(MessageContext)

  const editIncomingWebhook = useCallback(
    async (incomingWebhook: CreateOrUpdateIncomingWebhookRequest) => {
      try {
        await editIncomingWebhookMutation({
          variables: incomingWebhook,
        })
        showMessage(`incoming webhook edited: ${incomingWebhook.alias}`)
        history.goBack()
      } catch (err) {
        setErrorMessage(getGQLError(err))
      }
    },
    [editIncomingWebhookMutation, showMessage, history]
  )

  const handleOnSubmit = useCallback(
    (e: FormEvent | MouseEvent) => {
      e.preventDefault()
      if (!isValidForm(formState)) {
        setErrorMessage('Please fill out correctly the mandatory fields.')
        return
      }
      const { alias } = formState.values
      editIncomingWebhook({ id: data.id, alias })
    },
    [formState, data, editIncomingWebhook]
  )

  return (
    <>
      <header>
        <h1>Edit incoming webhook #{data.id}</h1>
      </header>
      <section>
        <IncomingWebhookHelp data={data} />
        {errorMessage != null && <ErrorPanel title="Unable to edit incoming webhook">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField label="Alias" {...text('alias')} error={formState.errors.alias} required autoFocus />
        </form>
      </section>
      <footer>
        <Button title="Back to integrations" as={Link} to="/settings/integrations">
          Cancel
        </Button>
        <Button title="Edit incoming webhook" onClick={handleOnSubmit} variant="primary">
          Update
        </Button>
      </footer>
    </>
  )
}
