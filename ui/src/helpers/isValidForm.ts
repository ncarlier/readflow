import { FormState } from 'react-use-form-state'

import { FormMountValidity } from '../hooks/useOnMountInputValidator'

export function isValidForm(formState: FormState<any>, onMountValidator: FormMountValidity<any>) {
  const validity = { ...onMountValidator.validity, ...formState.validity }
  Object.keys(formState.values).forEach(key => {
    if (!validity[key]) {
      return false
    }
  })
  return true
}
