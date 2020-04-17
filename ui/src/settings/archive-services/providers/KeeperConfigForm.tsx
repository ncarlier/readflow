import React from 'react'
import { useFormState } from 'react-use-form-state'

import FormInputField from '../../../components/FormInputField'

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
  // eslint-disable-next-line @typescript-eslint/camelcase
  api_key: ''
}

export default ({ onChange, config = defaultConfig }: Props) => {
  const [formState, { url, text }] = useFormState<KeeperConfigFormFields>(config, {
    onChange: (_e, _stateValues, nextStateValues) => onChange(nextStateValues)
  })

  return (
    <>
      <FormInputField label="Endpoint" {...url('endpoint')} error={formState.errors.endpoint} required />
      <FormInputField label="API key" {...text('api_key')} error={formState.errors.api_key} required />
    </>
  )
}
