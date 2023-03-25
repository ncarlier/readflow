import React from 'react'
import { useFormState } from 'react-use-form-state'

import { FormInputField } from '../../../../../components'

interface Config {
  endpoint: string
}

interface Secrets {
  api_key: string
}

type ConfigFormFields = Config & Secrets

export const marshal = (config: ConfigFormFields) : string[] => [
  JSON.stringify(config, ['endpoint']),
  JSON.stringify(config, ['api_key']),
]

interface Props {
  onChange(config: any): void
  config?: ConfigFormFields
}

const defaultConfig = {
  endpoint: 'https://api.nunux.org/keeper/v2/documents',
}

export const ConfigForm = ({ onChange, config }: Props) => {
  const [formState, { url, text }] = useFormState<ConfigFormFields>(
    { ...defaultConfig, ...config },
    {
      onChange: (_e, _stateValues, nextStateValues) => onChange(nextStateValues),
    }
  )

  return (
    <>
      <FormInputField label="Endpoint" {...url('endpoint')} error={formState.errors.endpoint} required />
      <FormInputField label="API key" {...text('api_key')} error={formState.errors.api_key} required />
    </>
  )
}
