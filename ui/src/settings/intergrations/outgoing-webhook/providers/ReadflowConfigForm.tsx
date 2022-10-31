import React from 'react'
import { useFormState } from 'react-use-form-state'

import { FormInputField } from '../../../../components'
import { API_BASE_URL } from '../../../../constants'

interface ReadflowConfigFormFields {
  endpoint: string
  api_key: string
}

interface Props {
  onChange(config: any): void
  config?: ReadflowConfigFormFields
}

const defaultConfig = {
  endpoint: API_BASE_URL
}

export const ReadflowConfigForm = ({ onChange, config }: Props) => {
  const [formState, { url, text }] = useFormState<ReadflowConfigFormFields>(
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
