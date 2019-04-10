import React, { useCallback, useState } from 'react'

import { useFormState } from 'react-use-form-state'
import { useMutation } from 'react-apollo-hooks'

import Button from '../../common/Button'
import FormInputField from '../../common/FormInputField'
import ErrorPanel from '../../error/ErrorPanel'
import { getGQLError, isValidForm } from '../../common/helpers'
import { History } from 'history'
import { updateCacheAfterUpdate } from './cache'
import { ArchiveService } from './models'
import { CreateOrUpdateArchiveService } from './queries'
import FormSelectField from '../../common/FormSelectField'
import FormCheckboxField from '../../common/FormCheckboxField'
import KeeperConfigForm from './providers/KeeperConfigForm'
import { IMessageDispatchProps, connectMessageDispatch } from '../../containers/MessageContainer'

interface EditArchiveServiceFormFields {
  alias: string
  provider: string
  is_default: boolean
}

type Props = {
  data: ArchiveService
  history: History
}

type AllProps = Props & IMessageDispatchProps

export const EditArchiveServiceForm = ({ data, history, showMessage }: AllProps) => {
  const [config, setConfig] = useState<any>(JSON.parse(data.config)) 
  const [errorMessage, setErrorMessage] = useState<string | null>(null) 
  const [formState, { text, select, checkbox }] = useFormState<EditArchiveServiceFormFields>({
    alias: data.alias,
    provider: data.provider,
    is_default: data.is_default
  })
  const editArchiveServiceMutation = useMutation<ArchiveService>(CreateOrUpdateArchiveService)

  const editArchiveService = async (service: ArchiveService) => {
    try{
      const res = await editArchiveServiceMutation({
        variables: service,
        update: updateCacheAfterUpdate,
      })
      showMessage(`Archive service edited: ${service.alias}`)
      history.goBack()
    } catch (err) {
      setErrorMessage(getGQLError(err))
    }
  }

  const handleOnSubmit = useCallback(() => {
    if (!isValidForm(formState) || !config) {
      setErrorMessage("Please fill out correctly the mandatory fields.")
      return
    }
    const { alias, provider, is_default } = formState.values
    editArchiveService({id: data.id, alias, provider, is_default, config: JSON.stringify(config)})
  }, [formState, config])

  return (
    <>
      <header>
        <h1>Edit archive service #{data.id}</h1>
      </header>
      <section>
        {errorMessage != null &&
          <ErrorPanel title="Unable to edit archive service">
            {errorMessage}
          </ErrorPanel>
        }
        <form onSubmit={handleOnSubmit}>
          <FormInputField label="Alias"
            {...text('alias')}
            error={!formState.validity.alias}
            required />
          <FormSelectField label="Provider" {...select('provider')}>
            <option value="keeper">Keeper</option>
            <option value="wallabag">Wallabag</option>
          </FormSelectField>
          { formState.values.provider === 'keeper' && <KeeperConfigForm onChange={setConfig} config={config} /> }
          <FormCheckboxField label="To use by default" {...checkbox('is_default')} />
        </form>
      </section>
      <footer>
        <Button title="Back to archive services" to="/settings/archive-services">
          Cancel
        </Button>
        <Button
          title="Edit archive service"
          onClick={handleOnSubmit}
          primary>
          Update
        </Button>
      </footer>
    </>
  )
}

export default connectMessageDispatch(EditArchiveServiceForm)
