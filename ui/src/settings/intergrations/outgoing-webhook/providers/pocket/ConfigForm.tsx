import React, { useEffect } from 'react'
import { useFormState } from 'react-use-form-state'

import { FormInputField } from '../../../../../components'
import { AccountLinker } from '../../../components'

interface Config {
  username: string
}

interface Secrets {
  access_token: string
}

type ConfigFormFields = Config & Secrets

export const marshal = (config: ConfigFormFields) : string[] => [
  JSON.stringify(config, ['username']),
  JSON.stringify(config, ['access_token']),
]

interface Props {
  onChange(config: any): void
  config?: ConfigFormFields
}

export const ConfigForm = ({ onChange, config }: Props) => {
  const [formState, { text, password }] = useFormState<ConfigFormFields>(config, {
    onChange: (_e, _stateValues, nextStateValues) => onChange(nextStateValues),
  })
  // Hack used to update form if config change
  useEffect(() => {
    if (config) {
      formState.setField('access_token', config.access_token)
      formState.setField('username', config.username)
    }
  }, [formState, config])

  if (formState.values.access_token === '') {
    return <AccountLinker provider="pocket" />
  }

  return (
    <>
      <FormInputField
        label="Access token"
        {...password('access_token')}
        error={formState.errors.access_token}
        required
      />
      <FormInputField label="Username" {...text('username')} error={formState.errors.username} required />
    </>
  )
}
