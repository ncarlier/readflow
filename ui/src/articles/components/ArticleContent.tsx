/* eslint-disable @typescript-eslint/no-var-requires */
import mousetrap from 'mousetrap'
import React, { useEffect, useRef } from 'react'

import { useLocalConfiguration } from '../../contexts/LocalConfigurationContext'
import { Article } from '../models'
import styles from './ArticleContent.module.css'
import readable from './readable'
import { getEffectiveTheme } from '../../helpers'

interface Props {
  article: Article
}

const getHTMLContent = (body: string, theme: string) => `
  <head>
    <base target="_blank">
    <meta charset="utf-8" />
    <style>
      ${readable.style}
    </style>
  </head>
  <body data-theme="${theme}">
    ${body}
    <script>
      ${readable.script}
    </script>
  </body>
`

type KeyMap = {
  [key: string]: string
}
const keyMap: KeyMap = {
  'Delete': 'del',
  'Insert': 'ins',
}

export const ArticleContent = ({ article }: Props) => {
  const contentRef = useRef<HTMLDivElement>(null)
  const { localConfiguration } = useLocalConfiguration()

  useEffect(() => {
    if (contentRef.current) {
      const ifrm = document.createElement('iframe')
      //console.log('render')
      contentRef.current.innerHTML = ''
      contentRef.current.appendChild(ifrm)
      const doc = ifrm.contentWindow ? ifrm.contentWindow.document : ifrm.contentDocument
      if (doc) {
        doc.open()
        const root = doc.createElement('html')
        root.innerHTML = getHTMLContent(article.html || article.text, getEffectiveTheme(localConfiguration.theme))
        doc.appendChild(root)
        // Keyboard events have to propagate outside the iframe
        doc.onkeydown = function (e: KeyboardEvent) {
          if (e.key in keyMap) {
            mousetrap.trigger(keyMap[e.key])
          } else {
            mousetrap.trigger(e.key.toLowerCase())
          }
        }
        doc.close()
      }
      ifrm.focus()
    }
  }, [article, useLocalConfiguration])

  return <article className={styles.content} ref={contentRef} />
}
