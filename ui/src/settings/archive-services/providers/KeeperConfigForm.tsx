import React from 'react'
import { useFormState } from 'react-use-form-state'

import FormInputField from '../../../components/FormInputField'
import { isValidInput } from '../../../helpers'
import useOnMountInputValidator from '../../../hooks/useOnMountInputValidator'

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
        ref={onMountValidator.bind}
      />
      <FormInputField
        label="API key"
        {...text('api_key')}
        error={!isValidInput(formState, onMountValidator, 'api_key')}
        required
        ref={onMountValidator.bind}
      />
    </>
  )
}
