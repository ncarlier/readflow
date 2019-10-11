import React from 'react'
import { CheckboxProps } from 'react-use-form-state'

import { classNames } from '../helpers'

interface Props {
  label: string
  required?: boolean
  readOnly?: boolean
  error?: boolean
}

type AllProps = Props & CheckboxProps<any>

export default (props: AllProps) => {
  const { error, label, ...rest } = props
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
