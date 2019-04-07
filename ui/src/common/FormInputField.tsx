import React, { ReactNode } from 'react'
import { BaseInputProps } from 'react-use-form-state'
import { classNames } from './helpers'

type Props = {
  label: string
  required?: boolean
  readOnly?: boolean
  error?: boolean
  children?: ReactNode
}

type AllProps = Props & BaseInputProps

export default (props: AllProps) => {
  const {error, label, children, ...rest} = props

  if (rest.type === "checkbox") {
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

  const className = classNames(
    'form-group',
    error ? 'has-error' : null
  )

  let $input: ReactNode
  switch (rest.type) {
    case 'select':
      $input = <select {...rest}>
        {children}
      </select>
      break
    case 'textarea':
      $input = <textarea {...rest}></textarea>
    default:
      $input = <input {...rest} />
      break;
  }

  return (
    <div className={className}>
      {$input}
      <label htmlFor={rest.name} className="control-label">{label}</label>
      <i className="bar"></i>
    </div>
  )
}
