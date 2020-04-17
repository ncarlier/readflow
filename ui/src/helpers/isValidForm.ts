import { FormState } from 'react-use-form-state'

export function isValidForm(formState: FormState<any>) {
  const { validity } = formState

  for (const key of Object.keys(validity)) {
    if (!validity[key]) {
      return false
    }
  }
  return true
}
