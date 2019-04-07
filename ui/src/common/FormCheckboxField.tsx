import React, { ReactNode } from 'react'
import { CheckboxProps } from 'react-use-form-state'
import { classNames } from './helpers'

type Props = {
  label: string
  required?: boolean
  readOnly?: boolean
  error?: boolean
}

type AllProps = Props & CheckboxProps

export default (props: AllProps) => {
  const {error, label, ...rest} = props
  const className = classNames(
    'checkbox',
     error ? 'has-error' : null 
  )
  return (
    <div className={className}>
      <label>
        <input {...rest} />
        <i className="helper"></i>{label}
      </label>
    </div>
  )
}
