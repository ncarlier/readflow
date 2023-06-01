import React, { forwardRef, ReactNode, Ref, useState } from 'react'
import Editor from 'react-simple-code-editor'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import http_lang from 'highlight.js/lib/languages/http'
import text_lang from 'highlight.js/lib/languages/plaintext'
import evalfilter_lang from './highlight/evalfiter'
import fast_lang from './highlight/fasttemplate'
hljs.registerLanguage('script', evalfilter_lang)
hljs.registerLanguage('headers', http_lang)
hljs.registerLanguage('template', fast_lang)
hljs.registerLanguage('text', text_lang)

import { BaseInputProps, Omit } from 'react-use-form-state'

import { classNames } from '../helpers'

interface Props {
  label: string
  language?: 'script' | 'headers' | 'template' | 'text'
  required?: boolean
  readOnly?: boolean
  pattern?: string
  maxLength?: number
  error?: string
  children?: ReactNode
}

type AllProps = Props & Omit<BaseInputProps<any>, 'type'>

export const FormCodeEditorField = forwardRef((props: AllProps, ref: Ref<any>) => {
  const { error, label, children, language = 'text', value, ...rest } = { ...props, ref }
  const [code, setCode] = useState(value)

  const className = classNames('form-group', error ? 'has-error' : null)

  return (
    <div className={className}>
      <label htmlFor={rest.name} className="control-label alt">
        {label}
      </label>
      <Editor
        {...rest}
        value={code}
        onValueChange={code => setCode(code)}
        highlight={code => hljs.highlight(code, {language}).value}
      />
      <i className="bar" />
      {!!error && <><span className="helper">{error}</span><br /></>}
      {children}
    </div>
  )
})
