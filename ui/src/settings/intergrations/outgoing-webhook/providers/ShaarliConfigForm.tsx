import React from 'react'
import { useFormState } from 'react-use-form-state'

import { FormCheckboxField, FormInputField } from '../../../../components'

interface ShaarliConfigFormFields {
  endpoint: string
  secret: string
  private: boolean
}

interface Props {
  onChange(config: any): void
  config?: ShaarliConfigFormFields
}

const defaultConfig = {
  endpoint: 'https://demo.shaarli.org',
}

export const ShaarliConfigForm = ({ onChange, config }: Props) => {
  const [formState, { url, checkbox, password }] = useFormState<ShaarliConfigFormFields>(
    { ...defaultConfig, ...config },
    {
      onChange: (_e, _stateValues, nextStateValues) => onChange(nextStateValues),
    }
  )

  return (
    <>
      <FormInputField label="Endpoint" {...url('endpoint')} error={formState.errors.endpoint} required />
      <FormInputField label="Secret" {...password('secret')} error={formState.errors.secret} required />
      <FormCheckboxField label="Public" {...checkbox('secret')} error={formState.errors.private} />
    </>
  )
}
