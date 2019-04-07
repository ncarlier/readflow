import React, { useCallback, useState } from 'react'

import { useFormState } from 'react-use-form-state'
import { useMutation } from 'react-apollo-hooks'
import { RouteComponentProps } from 'react-router'

import Panel from '../../common/Panel'
import Button from '../../common/Button'
import { usePageTitle } from '../../hooks'
import FormInputField from '../../common/FormInputField'
import { CreateOrUpdateApiKey } from './queries'
import ErrorPanel from '../../error/ErrorPanel'
import { getGQLError, isValidForm } from '../../common/helpers'
import { updateCacheAfterCreate } from './cache'
import { connectMessageDispatch, IMessageDispatchProps } from '../../containers/MessageContainer'

interface AddApiKeyFormFields {
  alias: string
}

type AllProps = RouteComponentProps<{}> & IMessageDispatchProps

export const AddApiKeyForm = ({history, showMessage }: AllProps) => {
  usePageTitle('Settings - Add new API key')

  const [errorMessage, setErrorMessage] = useState<string | null>(null) 
  const [formState, { text }] = useFormState<AddApiKeyFormFields>()
  const addApiKeyMutation = useMutation<AddApiKeyFormFields>(CreateOrUpdateApiKey)

  const addApiKey = async (apiKey: {alias: string | null}) => {
    try{
      const res = await addApiKeyMutation({
        variables: apiKey,
        update: updateCacheAfterCreate
      })
      showMessage(`New API key: ${res.data.createOrUpdateAPIKey.id}`)
      // console.log('New API key', res)
      history.goBack()
    } catch (err) {
      setErrorMessage(getGQLError(err))
    }
  }

  const handleOnClick = useCallback(() => {
    if (!isValidForm(formState)) {
      setErrorMessage("Please fill out correctly the mandatory fields.")
      return
    }
    const {alias}  = formState.values
    addApiKey({alias})
  }, [formState])

  return (
    <Panel>
      <header>
        <h1>Add new API key</h1>
        <Button title="Back to API keys" to="/settings/api-keys">
          Cancel
        </Button>
        <Button
          title="Add API key"
          onClick={handleOnClick}
          primary>
          Add
        </Button>
      </header>
      <section>
        {errorMessage != null &&
          <ErrorPanel title="Unable to add new API key">
            {errorMessage}
          </ErrorPanel>
        }
        <form>
          <FormInputField label="Alias"
            {...text('alias')}
            error={!formState.validity.alias}
            required />
        </form>
      </section>
    </Panel>
  )
}

export default connectMessageDispatch(AddApiKeyForm)
