import React, { forwardRef, Ref } from 'react'
import { CheckboxProps } from 'react-use-form-state'

import { classNames } from '../helpers'

interface Props {
  label: string
  required?: boolean
  readOnly?: boolean
  error?: string
}

type AllProps = Props & CheckboxProps<any>

export default forwardRef((props: AllProps, ref: Ref<any>) => {
  const { error, label, ...rest } = { ...props, ref }
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
})
