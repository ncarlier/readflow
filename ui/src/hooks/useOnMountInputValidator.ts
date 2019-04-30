import { useCallback, useEffect, useRef, useState } from 'react'
import { StateValidity } from 'react-use-form-state'

export interface FormMountValidity<T> {
  bind: (imput: any) => void
  validity: StateValidity<T>
}

export default <T>(validity: StateValidity<T>): FormMountValidity<T> => {
  const [onMountValidity, setOnMountValidity] = useState(validity)
  const inputRefs = useRef<any>({})
  const bind = useCallback(input => {
    if (input) {
      inputRefs.current[input.name] = input
    }
  }, [])
  useEffect(() => {
    for (const name in inputRefs.current) {
      setOnMountValidity(validity => ({
        ...validity,
        [name]: inputRefs.current[name].validity.valid
      }))
    }
  }, [])
  return {
    bind,
    validity: onMountValidity
  }
}
