/* eslint-disable @typescript-eslint/camelcase */
import { History } from 'history'
import React, { FormEvent, useCallback, useContext, useState } from 'react'
import { useMutation } from '@apollo/client'
import { Link } from 'react-router-dom'
import { useFormState } from 'react-use-form-state'

import Button from '../../../components/Button'
import FormCheckboxField from '../../../components/FormCheckboxField'
import FormInputField from '../../../components/FormInputField'
import FormSelectField from '../../../components/FormSelectField'
import { MessageContext } from '../../../context/MessageContext'
import ErrorPanel from '../../../error/ErrorPanel'
import { getGQLError, isValidForm } from '../../../helpers'
import { OutgoingWebhook, CreateOrUpdateOutgoingWebhookResponse } from './models'
import KeeperConfigForm from './providers/KeeperConfigForm'
import { CreateOrUpdateOutgoingWebhook } from './queries'
import WallabagConfigForm from './providers/WallabagConfigForm'
import GenericConfigForm from './providers/GenericConfigForm'

interface EditOutgoingWebhookFormFields {
  alias: string
  provider: string
  is_default: boolean
}

interface Props {
  data: OutgoingWebhook
  history: History
}

export default ({ data, history }: Props) => {
  const [config, setConfig] = useState<any>(JSON.parse(data.config))
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text, select, checkbox }] = useFormState<EditOutgoingWebhookFormFields>({
    alias: data.alias,
    provider: data.provider,
    is_default: data.is_default,
  })
  const [editOutgoingWebhookMutation] = useMutation<CreateOrUpdateOutgoingWebhookResponse, OutgoingWebhook>(
    CreateOrUpdateOutgoingWebhook
  )
  const { showMessage } = useContext(MessageContext)

  const editOutgoingWebhook = useCallback(
    async (webhook: OutgoingWebhook) => {
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
      editOutgoingWebhook({ id: data.id, alias, provider, is_default, config: JSON.stringify(config) })
    },
    [data, formState, config, editOutgoingWebhook]
  )

  return (
    <>
      <header>
        <h1>Edit outgoing webhook #{data.id}</h1>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to edit outgoing webhook">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField label="Alias" {...text('alias')} error={formState.errors.alias} required autoFocus />
          <FormSelectField label="Provider" {...select('provider')}>
            <option value="generic">Generic webhook</option>
            <option value="keeper">Keeper</option>
            <option value="wallabag">Wallabag</option>
          </FormSelectField>
          {formState.values.provider === 'generic' && <GenericConfigForm onChange={setConfig} />}
          {formState.values.provider === 'keeper' && <KeeperConfigForm onChange={setConfig} config={config} />}
          {formState.values.provider === 'wallabag' && <WallabagConfigForm onChange={setConfig} config={config} />}
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
