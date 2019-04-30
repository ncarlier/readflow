import { useCallback, useRef, useState } from 'react'
import { StateValidity } from 'react-use-form-state'

export default <T>(form: StateValidity<T>) => {
  const [validity, setValidity] = useState(form)
  const inputRefs = useRef<any>({})
  const bind = useCallback(input => {
    if (input) {
      inputRefs.current[input.name] = input
    }
  }, [])
  const validate = () => {
    for (const name in inputRefs.current) {
      setValidity(validityState => ({
        ...validityState,
        [name]: inputRefs.current[name].validity.valid
      }))
    }
  }
  return { bind, validate, validity }
}
