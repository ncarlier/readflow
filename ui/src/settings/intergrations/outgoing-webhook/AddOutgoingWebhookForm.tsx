import React, { FormEvent, useCallback, useEffect, useState } from 'react'
import { useMutation } from '@apollo/client'
import { RouteComponentProps } from 'react-router'
import { Link } from 'react-router-dom'
import { useFormState } from 'react-use-form-state'

import { useMessage } from '../../../contexts/MessageContext'
import { Button, ErrorPanel, FormCheckboxField, FormInputField, FormSelectField, Panel } from '../../../components'
import { getGQLError, isValidForm } from '../../../helpers'
import { usePageTitle } from '../../../hooks'
import { updateCacheAfterCreate } from './cache'
import { CreateOrUpdateOutgoingWebhookResponse, CreateOrUpdateOutgoingWebhookRequest, Provider } from './models'
import providers from './providers'
import { CreateOrUpdateOutgoingWebhook } from './queries'
import OutgoingWebhookHelp from './OutgoingWebhookHelp'

interface AddOutgoingWebhookFormFields {
  alias: string
  provider: Provider
  isDefault: boolean
}

const getFormStateFromQueryParams = (qs: string): AddOutgoingWebhookFormFields => {
  const params = new URLSearchParams(qs)
  return {
    alias: params.get('alias') || '',
    provider: (params.get('provider') as Provider) || '',
    isDefault: false,
  }
}

const getConfigFromQueryParams = (qs: string) => {
  const params = new URLSearchParams(qs)
  const result: any = {}
  params.forEach((v, k) => (result[k] = v))
  return result
}

export default ({ history, location: { search } }: RouteComponentProps) => {
  usePageTitle('Settings - Add new outgoing webhook')
  const [config, setConfig] = useState<any>(getConfigFromQueryParams(search))
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const { showMessage } = useMessage()
  const [formState, { text, checkbox, select }] = useFormState<AddOutgoingWebhookFormFields>(
    getFormStateFromQueryParams(search)
  )

  useEffect(() => {
    setConfig(getConfigFromQueryParams(search))
  }, [search])

  const [addOutgoingWebhookMutation] = useMutation<
    CreateOrUpdateOutgoingWebhookResponse,
    CreateOrUpdateOutgoingWebhookRequest
  >(CreateOrUpdateOutgoingWebhook)

  const addOutgoingWebhook = useCallback(
    async (outgoingWebhook: CreateOrUpdateOutgoingWebhookRequest) => {
      try {
        const res = await addOutgoingWebhookMutation({
          variables: outgoingWebhook,
          update: updateCacheAfterCreate,
        })
        if (res.data) {
          showMessage(`New outgoing webhook: ${res.data.createOrUpdateOutgoingWebhook.alias}`)
        }
        history.replace('/settings/integrations')
      } catch (err) {
        setErrorMessage(getGQLError(err))
      }
    },
    [addOutgoingWebhookMutation, showMessage, history]
  )

  const handleOnSubmit = useCallback(
    (e: FormEvent | MouseEvent) => {
      e.preventDefault()
      if (!isValidForm(formState) || !config) {
        console.log(formState, config)
        setErrorMessage('Please fill out correctly the mandatory fields.')
        return
      }
      const { alias, provider, isDefault } = formState.values
      const [conf, secrets] = providers[provider].marshal(config)
      addOutgoingWebhook({ alias, provider, is_default: isDefault, config: conf, secrets })
    },
    [formState, config, addOutgoingWebhook]
  )

  const { provider } = formState.values
  const ProviderConfigForm = provider ? providers[provider].form : null

  return (
    <Panel>
      <header>
        <h1>Add new outgoing webhook</h1>
      </header>
      <section>
        <OutgoingWebhookHelp />
        {errorMessage != null && <ErrorPanel title="Unable to add new outgoing webhook">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField label="Alias" {...text('alias')} error={formState.errors.alias} required pattern=".*\S+.*" maxLength={32} autoFocus />
          <FormSelectField label="Provider" {...select('provider')} error={formState.errors.provider} required>
            <option value="">Please select a webhook provider</option>
            {Object.entries(providers).map(([key, p]) => <option key={`provider-${key}`} value={key}>{p.label}</option>)}
          </FormSelectField>
          { ProviderConfigForm && <ProviderConfigForm onChange={setConfig} config={config} locked={false} /> }
          <FormCheckboxField label="To use by default" {...checkbox('isDefault')} />
        </form>
      </section>
      <footer>
        <Button title="Back to integrations" as={Link} to="/settings/integrations">
          Cancel
        </Button>
        <Button title="Add outgoing webhook" onClick={handleOnSubmit} variant="primary">
          Add
        </Button>
      </footer>
    </Panel>
  )
}
