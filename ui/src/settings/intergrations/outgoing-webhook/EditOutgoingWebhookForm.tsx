import { History } from 'history'
import React, { FormEvent, useCallback, useState } from 'react'
import { useMutation } from '@apollo/client'
import { Link } from 'react-router-dom'
import { useFormState } from 'react-use-form-state'

import { useMessage } from '../../../contexts'
import { Button, ErrorPanel, FormCheckboxField, FormInputField, FormSelectField } from '../../../components'
import { getGQLError, isValidForm } from '../../../helpers'
import {
  OutgoingWebhook,
  CreateOrUpdateOutgoingWebhookResponse,
  CreateOrUpdateOutgoingWebhookRequest,
  Provider,
} from './models'
import providers from './providers'
import { CreateOrUpdateOutgoingWebhook } from './queries'
import OutgoingWebhookHelp from './OutgoingWebhookHelp'

interface EditOutgoingWebhookFormFields {
  alias: string
  provider: Provider
  is_default: boolean
}

interface Props {
  data: OutgoingWebhook
  history: History
}

export default ({ data, history }: Props) => {
  const [config, setConfig] = useState<any>({
    ...(data.secrets.reduce((a, v) => ({ ...a, [v]: ''}), {})),
    ...JSON.parse(data.config),
  })
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text, select, checkbox }] = useFormState<EditOutgoingWebhookFormFields>({
    alias: data.alias,
    provider: data.provider,
    is_default: data.is_default,
  })
  const [editOutgoingWebhookMutation] = useMutation<
    CreateOrUpdateOutgoingWebhookResponse,
    CreateOrUpdateOutgoingWebhookRequest
  >(CreateOrUpdateOutgoingWebhook)
  const { showMessage } = useMessage()

  const editOutgoingWebhook = useCallback(
    async (webhook: CreateOrUpdateOutgoingWebhookRequest) => {
      try {
        await editOutgoingWebhookMutation({
          variables: webhook,
        })
        showMessage(`Outgoing webhook edited: ${webhook.alias}`)
        history.goBack()
      } catch (err) {
        setErrorMessage(getGQLError(err))
      }
    },
    [editOutgoingWebhookMutation, showMessage, history]
  )

  const handleOnSubmit = useCallback(
    (e: FormEvent | MouseEvent) => {
      e.preventDefault()
      if (!isValidForm(formState) || !config) {
        setErrorMessage('Please fill out correctly the mandatory fields.')
        return
      }
      const { alias, provider, is_default } = formState.values
      const [conf, secrets] = providers[provider].marshal(config)
      editOutgoingWebhook({ id: data.id, alias, provider, is_default, config: conf, secrets })
    },
    [data, formState, config, editOutgoingWebhook]
  )
  
  const ProviderConfigForm = providers[formState.values.provider].form

  return (
    <>
      <header>
        <h1>Edit outgoing webhook #{data.id}</h1>
      </header>
      <section>
        <OutgoingWebhookHelp />
        {errorMessage != null && <ErrorPanel title="Unable to edit outgoing webhook">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField label="Alias" {...text('alias')} error={formState.errors.alias} required pattern=".*\S+.*" maxLength={32} autoFocus />
          <FormSelectField label="Provider" {...select('provider')}>
            {Object.entries(providers).map(([key, p]) => <option key={`provider-${key}`} value={key}>{p.label}</option>)}
          </FormSelectField>
          <ProviderConfigForm onChange={setConfig} config={config} />
          <FormCheckboxField label="To use by default" {...checkbox('is_default')} />
        </form>
      </section>
      <footer>
        <Button title="Back to integrations" as={Link} to="/settings/integrations">
          Cancel
        </Button>
        <Button title="Edit outgoing webhook" onClick={handleOnSubmit} variant="primary">
          Update
        </Button>
      </footer>
    </>
  )
}
