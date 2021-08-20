import React from 'react'
import { useFormState } from 'react-use-form-state'

import FormInputField from '../../../../components/FormInputField'
import FormSelectField from '../../../../components/FormSelectField'

interface S3ConfigFormFields {
  endpoint: string
  access_key_id: string
  access_key_secret: string
  region: string
  bucket: string
  format: string
}

interface Props {
  onChange(config: any): void
  config?: S3ConfigFormFields
}

const defaultConfig = {
  region: 'eu-west-3',
  format: 'html',
}

const formats: Map<string, string> = new Map([
  ['HTML', 'html'],
  ['HTML with images', 'html-single'],
  ['ZIP with HTML and images', 'zip'],
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

export default ({ onChange, config }: Props) => {
  const [formState, { url, text, password, select }] = useFormState<S3ConfigFormFields>(
    { ...defaultConfig, ...config },
    {
      onChange: (_e, _stateValues, nextStateValues) => onChange(nextStateValues),
    }
  )

  return (
    <>
      <FormInputField label="Endpoint" {...url('endpoint')} error={formState.errors.endpoint} required />
      <FormInputField label="Access Key" {...text('access_key_id')} error={formState.errors.access_key_id} required />
      <FormInputField
        label="Secret Key"
        {...password('access_key_secret')}
        error={formState.errors.access_key_secret}
        required
      />
      <FormInputField label="Region" {...text('region')} error={formState.errors.region} required />
      <FormInputField label="Bucket" {...text('bucket')} error={formState.errors.bucket} required />
      <FormSelectField label="Format" {...select('format')} error={formState.errors.format} required>
        <Formats />
      </FormSelectField>
    </>
  )
}
