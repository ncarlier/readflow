import React from 'react'
import { useFormState } from 'react-use-form-state'

import FormInputField from '../../../components/FormInputField'
import FormSelectField from '../../../components/FormSelectField'
import FormTextareaField from '../../../components/FormTextareaField'
import HelpLink from '../../../components/HelpLink'
import { isValidInput } from '../../../helpers'
import useOnMountInputValidator from '../../../hooks/useOnMountInputValidator'

interface WebhookConfigFormFields {
  endpoint: string
  contentType: string
  format?: string
}

interface Props {
  onChange(config: any): void
  config?: WebhookConfigFormFields
}

const contentTypes: Map<string, string> = new Map([
  ['JSON', 'application/json; charset=utf-8'],
  ['Text', 'text/plain; charset=utf-8'],
  ['HTML', 'text/html; charset=utf-8']
])

const defaultConfig = {
  endpoint: '',
  contentType: contentTypes.get('json') || '',
  format: ''
}

const ContentTypes = () => (
  <>
    {Array.from(contentTypes.keys()).map(key => (
      <option key={`content-type-${key}`} value={contentTypes.get(key)}>
        {key}
      </option>
    ))}
  </>
)

export default ({ onChange, config = defaultConfig }: Props) => {
  const [formState, { url, select, textarea }] = useFormState<WebhookConfigFormFields>(config, {
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
      <FormSelectField
        label="Content Type"
        {...select('contentType')}
        error={!isValidInput(formState, onMountValidator, 'contentType')}
        required
        ref={onMountValidator.bind}
      >
        <ContentTypes />
      </FormSelectField>
      <FormTextareaField
        label="Format"
        {...textarea('format')}
        error={!isValidInput(formState, onMountValidator, 'format')}
        ref={onMountValidator.bind}
      >
        <HelpLink href="https://about.readflow.app/docs/en/read-flow/organize/rules/#syntax">
          View format syntax
        </HelpLink>
      </FormTextareaField>
    </>
  )
}
