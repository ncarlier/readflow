import React, { useCallback, useState } from 'react'

import { useFormState } from 'react-use-form-state'
import { useMutation } from 'react-apollo-hooks'
import { RouteComponentProps } from 'react-router'

import Panel from '../../common/Panel'
import Button from '../../common/Button'
import { usePageTitle } from '../../hooks'
import { CreateOrUpdateArchiveService } from './queries'
import ErrorPanel from '../../error/ErrorPanel'
import { getGQLError, isValidForm } from '../../common/helpers'
import { updateCacheAfterCreate } from './cache'
import FormInputField from '../../common/FormInputField'
import FormCheckboxField from '../../common/FormCheckboxField'
import FormSelectField from '../../common/FormSelectField'
import KeeperConfigForm from './providers/KeeperConfigForm'
import { ArchiveService } from './models'
import { IMessageDispatchProps, connectMessageDispatch } from '../../containers/MessageContainer'

interface AddArchiveServiceFormFields {
  alias: string
  provider: string
  is_default: boolean
}

type AllProps = RouteComponentProps<{}> & IMessageDispatchProps

export const AddArchiveServiceForm = ({history, showMessage}: AllProps) => {
  usePageTitle('Settings - Add new archive provider')

  const [config, setConfig] = useState<any>(null) 
  const [errorMessage, setErrorMessage] = useState<string | null>(null) 
  const [formState, { text, checkbox, select }] = useFormState<AddArchiveServiceFormFields>(
    {provider: '', alias: '', is_default: false}
  )
  const addArchiveServiceMutation = useMutation<ArchiveService>(CreateOrUpdateArchiveService)

  const addArchiveService = async (service: ArchiveService) => {
    try{
      const res = await addArchiveServiceMutation({
        variables: service,
        update: updateCacheAfterCreate
      })
      showMessage(`New archive service: ${res.data.createOrUpdateArchiver.id}`)
      history.goBack()
    } catch (err) {
      setErrorMessage(getGQLError(err))
    }
  }

  const handleOnClick = useCallback(() => {
    if (!isValidForm(formState) || !config) {
      setErrorMessage("Please fill out correctly the mandatory fields.")
      return
    }
    const { alias, provider, is_default } = formState.values
    addArchiveService({alias, provider, is_default, config: JSON.stringify(config)})
  }, [formState, config])

  return (
    <Panel>
      <header>
        <h1>Add new archive service</h1>
        <Button title="Back to archive serices" to="/settings/archive-services">
          Cancel
        </Button>
        <Button
          title="Add archive service"
          onClick={handleOnClick}
          primary>
          Add
        </Button>
      </header>
      <section>
        {errorMessage != null &&
          <ErrorPanel title="Unable to add new archive service">
            {errorMessage}
          </ErrorPanel>
        }
        <form>
          <FormInputField label="Alias"
            {...text('alias')}
            error={!formState.validity.alias}
            required />
          <FormSelectField label="Provider"
            {...select('provider')}
            error={!formState.validity.provider}
            required>
            <option>Please select an archive provider</option>
            <option value="keeper">Keeper</option>
            <option value="wallabag">Wallabag</option>
          </FormSelectField>
          { formState.values.provider === 'keeper' && <KeeperConfigForm onChange={setConfig} /> }
          <FormCheckboxField label="To use by default" {...checkbox('is_default')} />
        </form>
      </section>
    </Panel>
  )
}

export default connectMessageDispatch(AddArchiveServiceForm)
