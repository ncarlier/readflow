import React from 'react'
import { useFormState } from 'react-use-form-state'

import { FormInputField, FormSecretInputField, FormSelectField } from '../../../../../components'

interface Config {
  endpoint: string
  access_key_id: string
  region: string
  bucket: string
  format: string
}

interface Secrets {
  access_key_secret: string
}

type ConfigFormFields = Config & Secrets

export const marshal = (config: ConfigFormFields) : string[] => [
  JSON.stringify(config, ['endpoint', 'access_key_id', 'region', 'bucket', 'format']),
  JSON.stringify(config, ['access_key_secret']),
]

interface Props {
  onChange(config: any): void
  config?: ConfigFormFields
  locked?: boolean
}

const defaultConfig = {
  region: 'eu-west-3',
  format: 'html',
}

const formats: Map<string, string> = new Map([
  ['HTML file', 'html'],
  ['Single HTML file with images', 'html-single'],
  ['ZIP file with HTML and images', 'zip'],
  ['EPUB file', 'epub'],
  ['PDF file', 'pdf'],
  ['Markdown file', 'md'],
])

const Formats = () => (
  <>
    {Array.from(formats.keys()).map((key) => (
      <option key={`content-type-${key}`} value={formats.get(key)}>
        {key}
      </option>
    ))}
  </>
)

export const ConfigForm = ({ onChange, config, locked = true }: Props) => {
  const [formState, { url, text, select }] = useFormState<ConfigFormFields>(
    { ...defaultConfig, ...config },
    {
      onChange: (_e, _stateValues, nextStateValues) => onChange(nextStateValues),
    }
  )

  return (
    <>
      <FormInputField label="Endpoint" {...url('endpoint')} error={formState.errors.endpoint} required />
      <FormInputField label="Access Key" {...text('access_key_id')} error={formState.errors.access_key_id} required />
      <FormSecretInputField
        label="Secret Key"
        {...text('access_key_secret')}
        error={formState.errors.access_key_secret}
        required
        locked={locked}
      />
      <FormInputField label="Region" {...text('region')} error={formState.errors.region} required />
      <FormInputField label="Bucket" {...text('bucket')} error={formState.errors.bucket} required />
      <FormSelectField label="Format" {...select('format')} error={formState.errors.format} required>
        <Formats />
      </FormSelectField>
    </>
  )
}
