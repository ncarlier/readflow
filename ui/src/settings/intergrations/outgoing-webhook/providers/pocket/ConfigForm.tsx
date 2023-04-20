import React, { useEffect } from 'react'
import { useFormState } from 'react-use-form-state'

import { FormInputField, FormSecretInputField } from '../../../../../components'
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
  locked?: boolean
}

export const ConfigForm = ({ onChange, config, locked = true }: Props) => {
  const [formState, { text }] = useFormState<ConfigFormFields>(config, {
    onChange: (_e, _stateValues, nextStateValues) => onChange(nextStateValues),
  })
  // Hack used to update form if config change
  useEffect(() => {
    if (config) {
      formState.setField('access_token', config.access_token)
      formState.setField('username', config.username)
    }
  }, [formState, config])

  if (formState.values.username === '') {
    return <AccountLinker provider="pocket" />
  }

  return (
    <>
      <FormSecretInputField
        label="Access token"
        {...text('access_token')}
        error={formState.errors.access_token}
        required
        readOnly
        locked={locked}
      />
      <FormInputField label="Username" {...text('username')} error={formState.errors.username} required readOnly />
    </>
  )
}
