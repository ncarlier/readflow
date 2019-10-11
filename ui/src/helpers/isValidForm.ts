import { FormState } from 'react-use-form-state'

import { FormMountValidity } from '../hooks/useOnMountInputValidator'

export function isValidForm(formState: FormState<any>, onMountValidator: FormMountValidity<any>) {
  const validity = { ...onMountValidator.validity, ...formState.validity }
  for (const key of Object.keys(formState.values)) {
    if (!validity[key]) {
      return false
    }
  }
  return true
}

export function isValidInput(formState: FormState<any>, onMountValidator: FormMountValidity<any>, name: string) {
  if (typeof formState.validity[name] === 'undefined') {
    return onMountValidator.validity[name]
  }
  return formState.validity[name]
}
