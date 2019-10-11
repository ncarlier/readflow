import React, { FormEvent, useCallback, useContext, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { RouteComponentProps } from 'react-router'
import { useFormState } from 'react-use-form-state'

import Button from '../../components/Button'
import FormCheckboxField from '../../components/FormCheckboxField'
import FormInputField from '../../components/FormInputField'
import FormSelectField from '../../components/FormSelectField'
import Panel from '../../components/Panel'
import { MessageContext } from '../../context/MessageContext'
import ErrorPanel from '../../error/ErrorPanel'
import { getGQLError, isValidForm, isValidInput } from '../../helpers'
import { usePageTitle } from '../../hooks'
import useOnMountInputValidator from '../../hooks/useOnMountInputValidator'
import { updateCacheAfterCreate } from './cache'
import { ArchiveService, CreateOrUpdateArchiveServiceResponse } from './models'
import KeeperConfigForm from './providers/KeeperConfigForm'
import WebhookConfigForm from './providers/WebhookConfigForm'
import { CreateOrUpdateArchiveService } from './queries'
import { Link } from 'react-router-dom'

interface AddArchiveServiceFormFields {
  alias: string
  provider: string
  isDefault: boolean
}

type AllProps = RouteComponentProps<{}>

export default ({ history }: AllProps) => {
  usePageTitle('Settings - Add new archive provider')

  const [config, setConfig] = useState<any>(null)
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const { showMessage } = useContext(MessageContext)
  const [formState, { text, checkbox, select }] = useFormState<AddArchiveServiceFormFields>({
    provider: '',
    alias: '',
    isDefault: false
  })
  const onMountValidator = useOnMountInputValidator(formState.validity)

  const addArchiveServiceMutation = useMutation<CreateOrUpdateArchiveServiceResponse, ArchiveService>(
    CreateOrUpdateArchiveService
  )

  const addArchiveService = async (service: ArchiveService) => {
    try {
      const res = await addArchiveServiceMutation({
        variables: service,
        update: updateCacheAfterCreate
      })
      if (res.data) {
        showMessage(`New archive service: ${res.data.createOrUpdateArchiver.alias}`)
      }
      history.goBack()
    } catch (err) {
      setErrorMessage(getGQLError(err))
    }
  }

  const handleOnSubmit = useCallback(
    (e: FormEvent | MouseEvent) => {
      e.preventDefault()
      if (!isValidForm(formState, onMountValidator) || !config) {
        setErrorMessage('Please fill out correctly the mandatory fields.')
        return
      }
      const { alias, provider, isDefault } = formState.values
      // eslint-disable-next-line @typescript-eslint/camelcase
      addArchiveService({ alias, provider, is_default: isDefault, config: JSON.stringify(config) })
    },
    [formState, config]
  )

  return (
    <Panel>
      <header>
        <h1>Add new archive service</h1>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to add new archive service">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField
            label="Alias"
            {...text('alias')}
            error={!isValidInput(formState, onMountValidator, 'alias')}
            required
            autoFocus
            ref={onMountValidator.bind}
          />
          <FormSelectField
            label="Provider"
            {...select('provider')}
            error={!isValidInput(formState, onMountValidator, 'provider')}
            required
            ref={onMountValidator.bind}
          >
            <option>Please select an archive provider</option>
            <option value="keeper">Keeper</option>
            <option value="webhook">Webhook</option>
          </FormSelectField>
          {formState.values.provider === 'keeper' && <KeeperConfigForm onChange={setConfig} />}
          {formState.values.provider === 'webhook' && <WebhookConfigForm onChange={setConfig} />}
          <FormCheckboxField label="To use by default" {...checkbox('isDefault')} />
        </form>
      </section>
      <footer>
        <Button title="Back to archive services" as={Link} to="/settings/archive-services">
          Cancel
        </Button>
        <Button title="Add archive service" onClick={handleOnSubmit} variant="primary">
          Add
        </Button>
      </footer>
    </Panel>
  )
}
