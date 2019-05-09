import React from 'react'
import { useFormState } from 'react-use-form-state'

import FormInputField from '../../../common/FormInputField'

interface WebhookConfigFormFields {
  endpoint: string
}

interface Props {
  onChange(config: any): void
  config?: WebhookConfigFormFields
}

const defaultConfig = {
  endpoint: ''
}

export default ({ onChange, config = defaultConfig }: Props) => {
  const [formState, { url, text }] = useFormState<WebhookConfigFormFields>(config, {
    onChange: (_e, stateValues, nextStateValues) => onChange(nextStateValues)
  })

  return (
    <>
      <FormInputField label="Endpoint" {...url('endpoint')} error={!formState.validity.endpoint} required />
    </>
  )
}
