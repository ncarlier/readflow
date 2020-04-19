/* eslint-disable @typescript-eslint/camelcase */
import { History } from 'history'
import React, { FormEvent, useCallback, useContext, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { Link } from 'react-router-dom'
import { useFormState } from 'react-use-form-state'

import Button from '../../components/Button'
import FormCheckboxField from '../../components/FormCheckboxField'
import FormInputField from '../../components/FormInputField'
import FormSelectField from '../../components/FormSelectField'
import { MessageContext } from '../../context/MessageContext'
import ErrorPanel from '../../error/ErrorPanel'
import { getGQLError, isValidForm } from '../../helpers'
import { updateCacheAfterUpdate } from './cache'
import { ArchiveService, CreateOrUpdateArchiveServiceResponse } from './models'
import KeeperConfigForm from './providers/KeeperConfigForm'
import WebhookConfigForm from './providers/WebhookConfigForm'
import { CreateOrUpdateArchiveService } from './queries'

interface EditArchiveServiceFormFields {
  alias: string
  provider: string
  is_default: boolean
}

interface Props {
  data: ArchiveService
  history: History
}

export default ({ data, history }: Props) => {
  const [config, setConfig] = useState<any>(JSON.parse(data.config))
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text, select, checkbox }] = useFormState<EditArchiveServiceFormFields>({
    alias: data.alias,
    provider: data.provider,
    is_default: data.is_default,
  })
  const [editArchiveServiceMutation] = useMutation<CreateOrUpdateArchiveServiceResponse, ArchiveService>(
    CreateOrUpdateArchiveService
  )
  const { showMessage } = useContext(MessageContext)

  const editArchiveService = useCallback(
    async (service: ArchiveService) => {
      try {
        await editArchiveServiceMutation({
          variables: service,
          update: updateCacheAfterUpdate,
        })
        showMessage(`Archive service edited: ${service.alias}`)
        history.goBack()
      } catch (err) {
        setErrorMessage(getGQLError(err))
      }
    },
    [editArchiveServiceMutation, showMessage, history]
  )

  const handleOnSubmit = useCallback(
    (e: FormEvent | MouseEvent) => {
      e.preventDefault()
      if (!isValidForm(formState) || !config) {
        setErrorMessage('Please fill out correctly the mandatory fields.')
        return
      }
      const { alias, provider, is_default } = formState.values
      editArchiveService({ id: data.id, alias, provider, is_default, config: JSON.stringify(config) })
    },
    [data, formState, config, editArchiveService]
  )

  return (
    <>
      <header>
        <h1>Edit archive service #{data.id}</h1>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to edit archive service">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField label="Alias" {...text('alias')} error={formState.errors.alias} required autoFocus />
          <FormSelectField label="Provider" {...select('provider')}>
            <option value="keeper">Keeper</option>
            <option value="webhook">Webhook</option>
          </FormSelectField>
          {formState.values.provider === 'keeper' && <KeeperConfigForm onChange={setConfig} config={config} />}
          {formState.values.provider === 'webhook' && <WebhookConfigForm onChange={setConfig} config={config} />}
          <FormCheckboxField label="To use by default" {...checkbox('is_default')} />
        </form>
      </section>
      <footer>
        <Button title="Back to archive services" as={Link} to="/settings/archive-services">
          Cancel
        </Button>
        <Button title="Edit archive service" onClick={handleOnSubmit} variant="primary">
          Update
        </Button>
      </footer>
    </>
  )
}
