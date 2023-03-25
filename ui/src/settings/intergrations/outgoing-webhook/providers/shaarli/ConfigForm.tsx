import React from 'react'
import { useFormState } from 'react-use-form-state'

import { FormCheckboxField, FormInputField } from '../../../../../components'

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
}

const defaultConfig = {
  endpoint: 'https://demo.shaarli.org',
}

export const ConfigForm = ({ onChange, config }: Props) => {
  const [formState, { url, checkbox, password }] = useFormState<ConfigFormFields>(
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
