import React, { forwardRef, Ref, useCallback, useState } from 'react'
import { TypeLessInputProps } from 'react-use-form-state'

import { classNames } from '../helpers'
import { Icon } from './Icon'
import { Button } from './Button'

interface Props {
  label: string
  required?: boolean
  pattern?: string
  maxLength?: number
  readOnly?: boolean
  autoFocus?: boolean
  error?: string
  locked?: boolean
}

type AllProps = Props & TypeLessInputProps<any>

export const FormSecretInputField = forwardRef((props: AllProps, ref: Ref<any>) => {
  const { error, label, locked = true, ...rest } = { ...props, ref }
  const [isLocked, setLocked] = useState(locked)

  const toggleLock = useCallback((e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    setLocked(!isLocked)
    return e.preventDefault()
  }, [isLocked])

  const className = classNames('form-group', error ? 'has-error' : null)

  return (
    <div className={className}>
      { isLocked && <Button icon='lock' onClick={toggleLock}>Unlock field</Button> }
      { !isLocked && <input type="text" {...rest} /> }
      <label htmlFor={rest.name} className="control-label">
        <Icon name={isLocked ? 'lock' : 'lock_open'} />
        {label}
      </label>
      { !isLocked && <i className="bar" /> }
      {!!error && <span className="helper">{error}</span>}
    </div>
  )
})
