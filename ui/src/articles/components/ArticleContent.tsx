/* eslint-disable @typescript-eslint/no-var-requires */
import mousetrap from 'mousetrap'
import React, { useEffect, useRef, useState } from 'react'

import { useLocalConfiguration } from '../../contexts/LocalConfigurationContext'
import { Article } from '../models'
import styles from './ArticleContent.module.css'
import readable from './readable'

const getMql = () => window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)')

interface Props {
  article: Article
}

const getHTMLContent = (body: string, theme: string) => `
<!DOCTYPE html>
<html lang="en">
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
</html>
`

type KeyMap = {
  [key: string]: string
}
const keyMap: KeyMap = {
  'Delete': 'del',
  'Insert': 'ins',
}

export const ArticleContent = ({ article }: Props) => {
  const [alreadyRendered, setAlreadyRendered] = useState(false)
  const contentRef = useRef<HTMLDivElement>(null)
  const { localConfiguration } = useLocalConfiguration()
  let { theme } = localConfiguration
  if (theme === 'auto') {
    const mql = getMql()
    theme = mql && mql.matches ? 'dark' : 'light'
  }

  useEffect(() => {
    if (contentRef.current && !alreadyRendered) {
      const ifrm = document.createElement('iframe')
      // console.log('render')
      contentRef.current.innerHTML = ''
      contentRef.current.appendChild(ifrm)
      const doc = ifrm.contentWindow ? ifrm.contentWindow.document : ifrm.contentDocument
      if (doc) {
        doc.open()
        doc.write(getHTMLContent(article.html || article.text, theme))
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
      setAlreadyRendered(true)
    }
  }, [alreadyRendered, article, theme])

  return <article className={styles.content} ref={contentRef} />
}
