import { History } from 'history'
import React, { FormEvent, useCallback, useContext, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { Link } from 'react-router-dom'
import { useFormState } from 'react-use-form-state'

import Button from '../../components/Button'
import FormInputField from '../../components/FormInputField'
import { MessageContext } from '../../context/MessageContext'
import ErrorPanel from '../../error/ErrorPanel'
import { getGQLError, isValidForm } from '../../helpers'
import { ApiKey, CreateOrUpdateApiKeyRequest, CreateOrUpdateApiKeyResponse } from './models'
import { CreateOrUpdateApiKey } from './queries'

interface EditApiKeyFormFields {
  alias: string
}

interface Props {
  data: ApiKey
  history: History
}

export default ({ data, history }: Props) => {
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text }] = useFormState<EditApiKeyFormFields>({ alias: data.alias })
  const [editApiKeyMutation] = useMutation<CreateOrUpdateApiKeyResponse, CreateOrUpdateApiKeyRequest>(
    CreateOrUpdateApiKey
  )
  const { showMessage } = useContext(MessageContext)

  const editApiKey = useCallback(
    async (apiKey: CreateOrUpdateApiKeyRequest) => {
      try {
        await editApiKeyMutation({
          variables: apiKey,
        })
        showMessage(`API key edited: ${apiKey.alias}`)
        history.goBack()
      } catch (err) {
        setErrorMessage(getGQLError(err))
      }
    },
    [editApiKeyMutation, showMessage, history]
  )

  const handleOnSubmit = useCallback(
    (e: FormEvent | MouseEvent) => {
      e.preventDefault()
      if (!isValidForm(formState)) {
        setErrorMessage('Please fill out correctly the mandatory fields.')
        return
      }
      const { alias } = formState.values
      editApiKey({ id: data.id, alias })
    },
    [formState, data, editApiKey]
  )

  return (
    <>
      <header>
        <h1>Edit API key #{data.id}</h1>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to edit API key">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField label="Alias" {...text('alias')} error={formState.errors.alias} required autoFocus />
        </form>
      </section>
      <footer>
        <Button title="Back to API keys" as={Link} to="/settings/api-keys">
          Cancel
        </Button>
        <Button title="Edit API key" onClick={handleOnSubmit} variant="primary">
          Update
        </Button>
      </footer>
    </>
  )
}
