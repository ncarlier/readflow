import React, { forwardRef, ReactNode, Ref } from 'react'
import Editor from 'react-simple-code-editor'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import evalfilter from './highlight/language'
hljs.registerLanguage('evalfilter', evalfilter)

import { BaseInputProps, Omit } from 'react-use-form-state'

import { classNames } from '../helpers'

interface Props {
  label: string
  required?: boolean
  readOnly?: boolean
  pattern?: string
  maxLength?: number
  error?: string
  children?: ReactNode
}

type AllProps = Props & Omit<BaseInputProps<any>, 'type'>

export const FormCodeEditorField = forwardRef((props: AllProps, ref: Ref<any>) => {
  const { error, label, children, ...rest } = { ...props, ref }

  const className = classNames('form-group', error ? 'has-error' : null)

  return (
    <div className={className}>
      <label htmlFor={rest.name} className="control-label-alt">
        {label}
      </label>
      <Editor
        {...rest}
        onValueChange={code => rest.value = code}
        highlight={code => hljs.highlight(code, {language: 'evalfilter'}).value}
      />
      <i className="bar" />
      {!!error && <><span className="helper">{error}</span><br /></>}
      {children}
    </div>
  )
})
