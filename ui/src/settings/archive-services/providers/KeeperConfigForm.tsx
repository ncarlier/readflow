import React, { useCallback } from 'react'

import { useFormState, FormState } from 'react-use-form-state'

import FormInputField from '../../../common/FormInputField'

interface KeeperConfigFormFields {
  endpoint: string
  api_key: string
}

type Props = {
  onChange(config: any): void
  config?: KeeperConfigFormFields 
}

export default ({onChange, config}: Props) => {
  const [formState, { url, text }] = useFormState<KeeperConfigFormFields>(config ? config : {
    endpoint: 'https://api.nunux.org/keeper/v2/documents',
    api_key: ''
  }, {
    onChange: (e, stateValues, nextStateValues) => onChange(nextStateValues) }
  )
  
  return (
    <>
      <FormInputField label="Endpoint"
        {...url('endpoint')}
        error={!formState.validity.endpoint}
        required />
      <FormInputField label="API key"
        {...text('api_key')}
        error={!formState.validity.api_key}
        required />
    </>
  )
}
