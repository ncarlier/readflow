import React, { useCallback, useState } from 'react'

import { useFormState } from 'react-use-form-state'
import { useMutation } from 'react-apollo-hooks'

import Button from '../../common/Button'
import FormInputField from '../../common/FormInputField'
import { CreateOrUpdateApiKey } from './queries'
import { ApiKey } from './models'
import ErrorPanel from '../../error/ErrorPanel'
import { getGQLError, isValidForm } from '../../common/helpers'
import { History } from 'history'
import { updateCacheAfterUpdate } from './cache'
import { IMessageDispatchProps, connectMessageDispatch } from '../../containers/MessageContainer'

interface EditApiKeyFormFields {
  alias: string
}

type Props = {
  data: ApiKey
  history: History
}

type AllProps = Props & IMessageDispatchProps

export const EditApiKeyForm = ({ data, history, showMessage }: AllProps) => {
  const [errorMessage, setErrorMessage] = useState<string | null>(null) 
  const [formState, { text }] = useFormState<EditApiKeyFormFields>({alias: data.alias})
  const editApiKeyMutation = useMutation<ApiKey>(CreateOrUpdateApiKey)

  const editApiKey = async (apiKey: {id: number, alias: string}) => {
    try{
      const res = await editApiKeyMutation({
        variables: apiKey,
        update: updateCacheAfterUpdate,
      })
      showMessage(`API key edited: ${apiKey.alias}`)
      history.goBack()
    } catch (err) {
      setErrorMessage(getGQLError(err))
    }
  }

  const handleOnSubmit = useCallback(() => {
    if (!isValidForm(formState)) {
      setErrorMessage("Please fill out correctly the mandatory fields.")
      return
    }
    const {alias} = formState.values
    editApiKey({id: data.id, alias})
  }, [formState])

  return (
    <>
      <header>
        <h1>Edit API key #{data.id}</h1>
      </header>
      <section>
        {errorMessage != null &&
          <ErrorPanel title="Unable to edit API key">
            {errorMessage}
          </ErrorPanel>
        }
        <form onSubmit={handleOnSubmit}>
          <FormInputField label="Alias"
            {...text('alias')}
            error={!formState.validity.alias}
            required />
        </form>
      </section>
      <footer>
        <Button title="Back to API keys" to="/settings/api-keys">
          Cancel
        </Button>
        <Button
          title="Edit API key"
          onClick={handleOnSubmit}
          primary>
          Update
        </Button>
      </footer>
    </>
  )
}

export default connectMessageDispatch(EditApiKeyForm)
