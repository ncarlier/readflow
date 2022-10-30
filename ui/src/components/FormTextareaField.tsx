import React, { forwardRef, ReactNode, Ref } from 'react'
import { BaseInputProps, Omit } from 'react-use-form-state'

import { classNames } from '../helpers'

interface Props {
  label: string
  required?: boolean
  readOnly?: boolean
  maxlength?: string
  error?: string
  children?: ReactNode
}

type AllProps = Props & Omit<BaseInputProps<any>, 'type'>

export const FormTextareaField = forwardRef((props: AllProps, ref: Ref<any>) => {
  const { error, label, children, ...rest } = { ...props, ref }

  const className = classNames('form-group', error ? 'has-error' : null)

  return (
    <div className={className}>
      <textarea {...rest} />
      <label htmlFor={rest.name} className="control-label">
        {label}
      </label>
      <i className="bar" />
      {children}
    </div>
  )
})
