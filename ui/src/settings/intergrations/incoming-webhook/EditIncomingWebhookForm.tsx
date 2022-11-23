import { History } from 'history'
import React, { FormEvent, useCallback, useState } from 'react'
import { useMutation } from '@apollo/client'
import { Link } from 'react-router-dom'
import { useFormState } from 'react-use-form-state'

import { useMessage } from '../../../contexts'
import { Button, ErrorPanel, FormCodeEditorField, FormInputField, HelpLink } from '../../../components'
import { getGQLError, isValidForm } from '../../../helpers'
import IncomingWebhookHelp from './IncomingWebhookHelp'
import { IncomingWebhook, CreateOrUpdateIncomingWebhookRequest, CreateOrUpdateIncomingWebhookResponse } from './models'
import { CreateOrUpdateIncomingWebhook } from './queries'

interface EditIncomingWebhookFormFields {
  alias: string
  script: string
}

interface Props {
  data: IncomingWebhook
  history: History
}

export default ({ data, history }: Props) => {
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text, textarea }] = useFormState<EditIncomingWebhookFormFields>(data)
  const [editIncomingWebhookMutation] = useMutation<
    CreateOrUpdateIncomingWebhookResponse,
    CreateOrUpdateIncomingWebhookRequest
  >(CreateOrUpdateIncomingWebhook)
  const { showMessage } = useMessage()

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
      editIncomingWebhook({ id: data.id, ...formState.values })
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
          <FormInputField label="Alias" {...text('alias')} error={formState.errors.alias} required pattern=".*\S+.*" maxLength={32} autoFocus />
          <FormCodeEditorField label="Script" language='script' {...textarea('script')} error={formState.errors.script} required pattern=".*\S+.*" maxLength={1024} >
            <HelpLink href="https://docs.readflow.app/en/integrations/incoming-webhook/scripting/">View script syntax</HelpLink>
          </FormCodeEditorField>
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
