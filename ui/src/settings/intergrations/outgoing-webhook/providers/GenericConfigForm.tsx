import React from 'react'
import { useFormState } from 'react-use-form-state'

import { FormInputField, FormSelectField, FormTextareaField, HelpLink } from '../../../../components'

interface GenericConfigFormFields {
  endpoint: string
  contentType: string
  format?: string
}

interface Props {
  onChange(config: any): void
  config?: GenericConfigFormFields
}

const contentTypes: Map<string, string> = new Map([
  ['JSON', 'application/json; charset=utf-8'],
  ['Text', 'text/plain; charset=utf-8'],
  ['HTML', 'text/html; charset=utf-8'],
])

const ContentTypes = () => (
  <>
    {Array.from(contentTypes.keys()).map((key) => (
      <option key={`content-type-${key}`} value={contentTypes.get(key)}>
        {key}
      </option>
    ))}
  </>
)

export const GenericConfigForm = ({ onChange, config }: Props) => {
  const [formState, { url, select, textarea }] = useFormState<GenericConfigFormFields>(config, {
    onChange: (_e, _stateValues, nextStateValues) => onChange(nextStateValues),
  })

  return (
    <>
      <FormInputField label="Endpoint" {...url('endpoint')} error={formState.errors.endpoint} required />
      <FormSelectField label="Content Type" {...select('contentType')} error={formState.errors.contentType} required>
        <ContentTypes />
      </FormSelectField>
      <FormTextareaField label="Format" {...textarea('format')} error={formState.errors.format}>
        <HelpLink href="https://docs.readflow.app/en/third-party/archive/webhook/#format">View format syntax</HelpLink>
      </FormTextareaField>
    </>
  )
}
