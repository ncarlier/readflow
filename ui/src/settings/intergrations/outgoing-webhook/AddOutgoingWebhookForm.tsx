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
import { KeeperConfigForm, GenericConfigForm, PocketConfigForm, S3ConfigForm, ShaarliConfigForm, WallabagConfigForm } from './providers'
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
      addOutgoingWebhook({ alias, provider, is_default: isDefault, config: JSON.stringify(config) })
    },
    [formState, config, addOutgoingWebhook]
  )

  return (
    <Panel>
      <header>
        <h1>Add new outgoing webhook</h1>
      </header>
      <section>
        <OutgoingWebhookHelp />
        {errorMessage != null && <ErrorPanel title="Unable to add new outgoing webhook">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField label="Alias" {...text('alias')} error={formState.errors.alias} required autoFocus />
          <FormSelectField label="Provider" {...select('provider')} error={formState.errors.provider} required>
            <option>Please select a webhook provider</option>
            <option value="generic">Generic webhook</option>
            <option value="keeper">Keeper</option>
            <option value="pocket">Pocket</option>
            <option value="s3">S3</option>
            <option value="shaarli">Shaarli</option>
            <option value="wallabag">Wallabag</option>
          </FormSelectField>
          {formState.values.provider === 'generic' && <GenericConfigForm onChange={setConfig} config={config} />}
          {formState.values.provider === 'keeper' && <KeeperConfigForm onChange={setConfig} config={config} />}
          {formState.values.provider === 'pocket' && <PocketConfigForm onChange={setConfig} config={config} />}
          {formState.values.provider === 's3' && <S3ConfigForm onChange={setConfig} config={config} />}
          {formState.values.provider === 'shaarli' && <ShaarliConfigForm onChange={setConfig} config={config} />}
          {formState.values.provider === 'wallabag' && <WallabagConfigForm onChange={setConfig} config={config} />}
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
