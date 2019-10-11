import React from 'react'
import { useFormState } from 'react-use-form-state'

import FormInputField from '../../../components/FormInputField'
import { isValidInput } from '../../../helpers'
import useOnMountInputValidator from '../../../hooks/useOnMountInputValidator'

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
  const [formState, { url }] = useFormState<WebhookConfigFormFields>(config, {
    onChange: (_e, stateValues, nextStateValues) => onChange(nextStateValues)
  })
  const onMountValidator = useOnMountInputValidator(formState.validity)

  return (
    <>
      <FormInputField
        label="Endpoint"
        {...url('endpoint')}
        error={!isValidInput(formState, onMountValidator, 'endpoint')}
        required
      />
    </>
  )
}
