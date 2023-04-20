import React from 'react'
import { useFormState } from 'react-use-form-state'

import { FormInputField, FormSecretInputField } from '../../../../../components'

interface Config {
  endpoint: string
  client_id: string
  username: string
}

interface Secrets {
  client_secret: string
  password: string
}

type ConfigFormFields = Config & Secrets

export const marshal = (config: ConfigFormFields) : string[] => [
  JSON.stringify(config, ['endpoint', 'client_id', 'username']),
  JSON.stringify(config, ['client_secret', 'password']),
]

interface Props {
  onChange(config: any): void
  config?: ConfigFormFields
  locked?: boolean
}

const defaultConfig = {
  endpoint: 'https://app.wallabag.it',
}

export const ConfigForm = ({ onChange, config, locked = true }: Props) => {
  const [formState, { url, text, password }] = useFormState<ConfigFormFields>(
    { ...defaultConfig, ...config },
    {
      onChange: (_e, _stateValues, nextStateValues) => onChange(nextStateValues),
    }
  )

  return (
    <>
      <FormInputField label="Endpoint" {...url('endpoint')} error={formState.errors.endpoint} required />
      <FormInputField label="Client ID" {...text('client_id')} error={formState.errors.client_id} required />
      <FormSecretInputField
        label="Client Secret"
        {...text('client_secret')}
        error={formState.errors.client_secret}
        required
        locked={locked}
      />
      <FormInputField label="Username" {...text('username')} error={formState.errors.username} required />
      <FormSecretInputField
        label="Password"
        {...password('password')}
        error={formState.errors.password}
        required
        locked={locked}
      />
    </>
  )
}
