import { History } from 'history'
import React, { FormEvent, useCallback, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { useFormState } from 'react-use-form-state'

import Button from '../../common/Button'
import FormInputField from '../../common/FormInputField'
import { getGQLError, isValidForm } from '../../common/helpers'
import { connectMessageDispatch, IMessageDispatchProps } from '../../containers/MessageContainer'
import ErrorPanel from '../../error/ErrorPanel'
import useOnMountInputValidator from '../../hooks/useOnMountInputValidator'
import { updateCacheAfterUpdate } from './cache'
import { ApiKey } from './models'
import { CreateOrUpdateApiKey } from './queries'

interface EditApiKeyFormFields {
  alias: string
}

interface Props {
  data: ApiKey
  history: History
}

type AllProps = Props & IMessageDispatchProps

export const EditApiKeyForm = ({ data, history, showMessage }: AllProps) => {
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text }] = useFormState<EditApiKeyFormFields>({ alias: data.alias })
  const onMountValidator = useOnMountInputValidator(formState.validity)
  const editApiKeyMutation = useMutation<ApiKey>(CreateOrUpdateApiKey)

  const editApiKey = async (apiKey: { id: number; alias: string }) => {
    try {
      await editApiKeyMutation({
        variables: apiKey,
        update: updateCacheAfterUpdate
      })
      showMessage(`API key edited: ${apiKey.alias}`)
      history.goBack()
    } catch (err) {
      setErrorMessage(getGQLError(err))
    }
  }

  const handleOnSubmit = useCallback(
    (e: FormEvent<HTMLFormElement>) => {
      e.preventDefault()
      if (!isValidForm(formState, onMountValidator)) {
        setErrorMessage('Please fill out correctly the mandatory fields.')
        return
      }
      const { alias } = formState.values
      editApiKey({ id: data.id, alias })
    },
    [formState]
  )

  return (
    <>
      <header>
        <h1>Edit API key #{data.id}</h1>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to edit API key">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField
            label="Alias"
            {...text('alias')}
            error={!formState.validity.alias}
            required
            ref={onMountValidator.bind}
          />
        </form>
      </section>
      <footer>
        <Button title="Back to API keys" to="/settings/api-keys">
          Cancel
        </Button>
        <Button title="Edit API key" onClick={handleOnSubmit} primary>
          Update
        </Button>
      </footer>
    </>
  )
}

export default connectMessageDispatch(EditApiKeyForm)
