import React, { forwardRef, ReactNode, Ref } from 'react'
import { BaseInputProps } from 'react-use-form-state'

import { classNames } from '../helpers'

interface Props {
  label: string
  required?: boolean
  readOnly?: boolean
  autoFocus?: boolean
  error?: string
  children?: ReactNode
}

type AllProps = Props & BaseInputProps<any>

export const FormInputField = forwardRef((props: AllProps, ref: Ref<any>) => {
  const { error, label, children, ...rest } = { ...props, ref }

  if (rest.type === 'checkbox') {
    const className = classNames('checkbox', error ? 'has-error' : null)
    return (
      <div className={className}>
        <label>
          <input {...rest} />
          <i className="helper" />
          {label}
        </label>
      </div>
    )
  }

  const className = classNames('form-group', error ? 'has-error' : null)

  let $input: ReactNode
  switch (rest.type) {
    case 'select':
      $input = <select {...rest}>{children}</select>
      break
    case 'textarea':
      $input = <textarea {...rest} />
      break
    default:
      $input = <input {...rest} />
  }

  return (
    <div className={className}>
      {$input}
      <label htmlFor={rest.name} className="control-label">
        {label}
      </label>
      <i className="bar" />
      {!!error && <span className="helper">{error}</span>}
    </div>
  )
})
