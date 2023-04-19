import React from 'react'
import { useFormState } from 'react-use-form-state'

import { FormCodeEditorField, FormInputField, HelpLink } from '../../../../../components'

interface Config {
  endpoint: string
  headers: Record<string, string>
  body: string
}

interface ConfigFormFields {
  endpoint: string
  headers: string
  body: string
}

export const marshal = (config: Config) : string[] => [JSON.stringify(config), '{}']

interface Props {
  onChange(config: any): void
  config: Config
}

const defautHeaders = {
  'Content-Type': 'application/json',
}

const HEADERS_REGEX = /^(\w+(?:-\w+)*: .+\n?)+$/

const validateHeaders = (value: string) => value.trim() && !HEADERS_REGEX.test(value) ? 'Invalid headers definition' : true

export const ConfigForm = ({ onChange, config: {endpoint, body, headers = defautHeaders} }: Props) => {
  const [formState, { url, textarea }] = useFormState<ConfigFormFields>({
    endpoint,
    body,
    headers: Object.keys(headers).map(k => `${k}: ${headers[k]}`).join('\n'),
  }, {
    onChange: (_e, _stateValues, nextStateValues) => onChange({
      ...nextStateValues,
      headers: nextStateValues.headers.split('\n').reduce(
        (obj: Record<string, string>, line) => {
          const [k, v] = line.split(': ')
          obj[k] = v
          return obj
        },
        {}
      ),
    }),
  })

  return (
    <>
      <FormInputField label="Endpoint" {...url('endpoint')} error={formState.errors.endpoint} required />
      <FormCodeEditorField label="HTTP headers" language="headers" {...textarea({name: 'headers', validate: validateHeaders})} error={formState.errors.headers} maxLength={1024} >
      </FormCodeEditorField>
      <FormCodeEditorField label="HTTP body" language="template" {...textarea('body')} error={formState.errors.body} maxLength={1024} >
        <HelpLink href="https://docs.readflow.app/en/integrations/outgoing-webhook/generic/#templating">View template syntax</HelpLink>
      </FormCodeEditorField>
    </>
  )
}
