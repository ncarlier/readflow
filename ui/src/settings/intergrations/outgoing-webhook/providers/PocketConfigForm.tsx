import React, { useEffect } from 'react'
import { useFormState } from 'react-use-form-state'

import FormInputField from '../../../../components/FormInputField'
import AccountLinker from '../../components/AccountLinker'

interface PocketConfigFormFields {
  access_token: string
  username: string
}

interface Props {
  onChange(config: any): void
  config: PocketConfigFormFields
}

export default ({ onChange, config }: Props) => {
  const [formState, { text, password }] = useFormState<PocketConfigFormFields>(config, {
    onChange: (_e, _stateValues, nextStateValues) => onChange(nextStateValues),
  })
  // Hack used to update form if config change
  useEffect(() => {
    formState.setField('access_token', config.access_token)
    formState.setField('username', config.username)
  }, [formState, config])

  if (formState.values.access_token === '') {
    return <AccountLinker provider="pocket" />
  }

  return (
    <>
      <FormInputField
        label="Access token"
        {...password('access_token')}
        error={formState.errors.access_token}
        required
      />
      <FormInputField label="Username" {...text('username')} error={formState.errors.username} required />
    </>
  )
}
