import React from 'react'
import { useFormState } from 'react-use-form-state'

import FormInputField from '../../../../components/FormInputField'

interface WallabagConfigFormFields {
  endpoint: string
  client_id: string
  client_secret: string
  username: string
  password: string
}

interface Props {
  onChange(config: any): void
  config?: WallabagConfigFormFields
}

const defaultConfig = {
  endpoint: 'https://app.wallabag.it',
}

export default ({ onChange, config }: Props) => {
  const [formState, { url, text, password }] = useFormState<WallabagConfigFormFields>(
    { ...defaultConfig, ...config },
    {
      onChange: (_e, _stateValues, nextStateValues) => onChange(nextStateValues),
    }
  )

  return (
    <>
      <FormInputField label="Endpoint" {...url('endpoint')} error={formState.errors.endpoint} required />
      <FormInputField label="Client ID" {...text('client_id')} error={formState.errors.client_id} required />
      <FormInputField
        label="Client Secret"
        {...text('client_secret')}
        error={formState.errors.client_secret}
        required
      />
      <FormInputField label="Username" {...text('username')} error={formState.errors.username} required />
      <FormInputField label="Password" {...password('password')} error={formState.errors.password} required />
    </>
  )
}
