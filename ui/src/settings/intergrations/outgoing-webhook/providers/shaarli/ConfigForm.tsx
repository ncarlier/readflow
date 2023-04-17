import React from 'react'
import { useFormState } from 'react-use-form-state'

import { FormCheckboxField, FormInputField, FormSecretInputField } from '../../../../../components'

interface Config {
  endpoint: string
  private: boolean
}

interface Secrets {
  secret: string
}

type ConfigFormFields = Config & Secrets

export const marshal = (config: ConfigFormFields) : string[] => [
  JSON.stringify(config, ['endpoint', 'private']),
  JSON.stringify(config, ['secret']),
]

interface Props {
  onChange(config: any): void
  config?: ConfigFormFields
  locked?: boolean
}

const defaultConfig = {
  endpoint: 'https://demo.shaarli.org',
}

export const ConfigForm = ({ onChange, config, locked = true }: Props) => {
  const [formState, { url, checkbox, password }] = useFormState<ConfigFormFields>(
    { ...defaultConfig, ...config },
    {
      onChange: (_e, _stateValues, nextStateValues) => onChange(nextStateValues),
    }
  )

  return (
    <>
      <FormInputField label="Endpoint" {...url('endpoint')} error={formState.errors.endpoint} required />
      <FormSecretInputField label="Secret" {...password('secret')} error={formState.errors.secret} required locked={locked} />
      <FormCheckboxField label="Public" {...checkbox('secret')} error={formState.errors.private} />
    </>
  )
}
