import React, { FormEvent, useCallback, useContext, useState } from 'react'
import { useMutation } from '@apollo/client'
import { RouteComponentProps } from 'react-router'
import { Link } from 'react-router-dom'
import { useFormState } from 'react-use-form-state'

import Button from '../../../components/Button'
import FormCheckboxField from '../../../components/FormCheckboxField'
import FormInputField from '../../../components/FormInputField'
import FormSelectField from '../../../components/FormSelectField'
import Panel from '../../../components/Panel'
import { MessageContext } from '../../../context/MessageContext'
import ErrorPanel from '../../../error/ErrorPanel'
import { getGQLError, isValidForm } from '../../../helpers'
import { usePageTitle } from '../../../hooks'
import { updateCacheAfterCreate } from './cache'
import { CreateOrUpdateOutgoingWebhookResponse, CreateOrUpdateOutgoingWebhookRequest } from './models'
import KeeperConfigForm from './providers/KeeperConfigForm'
import { CreateOrUpdateOutgoingWebhook } from './queries'
import WallabagConfigForm from './providers/WallabagConfigForm'
import GenericConfigForm from './providers/GenericConfigForm'

interface AddOutgoingWebhookFormFields {
  alias: string
  provider: string
  isDefault: boolean
}

export default ({ history }: RouteComponentProps) => {
  usePageTitle('Settings - Add new outgoing webhook')

  const [config, setConfig] = useState<any>(null)
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const { showMessage } = useContext(MessageContext)
  const [formState, { text, checkbox, select }] = useFormState<AddOutgoingWebhookFormFields>({
    provider: '',
    alias: '',
    isDefault: false,
  })

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
        history.goBack()
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
        setErrorMessage('Please fill out correctly the mandatory fields.')
        return
      }
      const { alias, provider, isDefault } = formState.values
      // eslint-disable-next-line @typescript-eslint/camelcase
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
        {errorMessage != null && <ErrorPanel title="Unable to add new outgoing webhook">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField label="Alias" {...text('alias')} error={formState.errors.alias} required autoFocus />
          <FormSelectField label="Provider" {...select('provider')} error={formState.errors.provider} required>
            <option>Please select an archive provider</option>
            <option value="generic">Generic webhook</option>
            <option value="keeper">Keeper</option>
            <option value="wallabag">Wallabag</option>
          </FormSelectField>
          {formState.values.provider === 'generic' && <GenericConfigForm onChange={setConfig} />}
          {formState.values.provider === 'keeper' && <KeeperConfigForm onChange={setConfig} />}
          {formState.values.provider === 'wallabag' && <WallabagConfigForm onChange={setConfig} />}
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
