import React from 'react'
import { useFormState } from 'react-use-form-state'

import FormInputField from '../../../../components/FormInputField'

interface KeeperConfigFormFields {
  endpoint: string
  api_key: string
}

interface Props {
  onChange(config: any): void
  config?: KeeperConfigFormFields
}

const defaultConfig = {
  endpoint: 'https://api.nunux.org/keeper/v2/documents',
}

export default ({ onChange, config }: Props) => {
  const [formState, { url, text }] = useFormState<KeeperConfigFormFields>(
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
