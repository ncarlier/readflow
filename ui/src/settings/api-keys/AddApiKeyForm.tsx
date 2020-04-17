import React, { FormEvent, useCallback, useContext, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'
import { RouteComponentProps } from 'react-router'
import { Link } from 'react-router-dom'
import { useFormState } from 'react-use-form-state'

import Button from '../../components/Button'
import FormInputField from '../../components/FormInputField'
import Panel from '../../components/Panel'
import { MessageContext } from '../../context/MessageContext'
import ErrorPanel from '../../error/ErrorPanel'
import { getGQLError, isValidForm } from '../../helpers'
import { usePageTitle } from '../../hooks'
import { updateCacheAfterCreate } from './cache'
import { CreateOrUpdateApiKeyRequest, CreateOrUpdateApiKeyResponse } from './models'
import { CreateOrUpdateApiKey } from './queries'

interface AddApiKeyFormFields {
  alias: string
}

type AllProps = RouteComponentProps<{}>

export default ({ history }: AllProps) => {
  usePageTitle('Settings - Add new API key')

  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [formState, { text }] = useFormState<AddApiKeyFormFields>()
  const [addApiKeyMutation] = useMutation<CreateOrUpdateApiKeyResponse, CreateOrUpdateApiKeyRequest>(CreateOrUpdateApiKey)
  const { showMessage } = useContext(MessageContext)

  const addApiKey = useCallback(
    async (apiKey: CreateOrUpdateApiKeyRequest) => {
      try {
        const res = await addApiKeyMutation({
          variables: apiKey,
          update: updateCacheAfterCreate
        })
        if (res.data) {
          showMessage(`New API key: ${res.data.createOrUpdateAPIKey.alias}`)
          // console.log('New API key', res)
          history.goBack()
        }
      } catch (err) {
        setErrorMessage(getGQLError(err))
      }
    },
    [addApiKeyMutation, showMessage, history]
  )

  const handleOnSubmit = useCallback(
    (e: FormEvent | MouseEvent) => {
      e.preventDefault()
      if (!isValidForm(formState)) {
        setErrorMessage('Please fill out correctly the mandatory fields.')
        return
      }
      const { alias } = formState.values
      addApiKey({ alias })
    },
    [formState, addApiKey]
  )

  return (
    <Panel>
      <header>
        <h1>Add new API key</h1>
      </header>
      <section>
        {errorMessage != null && <ErrorPanel title="Unable to add new API key">{errorMessage}</ErrorPanel>}
        <form onSubmit={handleOnSubmit}>
          <FormInputField label="Alias" {...text('alias')} error={formState.errors.alias} required autoFocus />
        </form>
      </section>
      <footer>
        <Button title="Back to API keys" as={Link} to="/settings/api-keys">
          Cancel
        </Button>
        <Button title="Add API key" onClick={handleOnSubmit} variant="primary">
          Add
        </Button>
      </footer>
    </Panel>
  )
}
