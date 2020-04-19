/* eslint-disable @typescript-eslint/no-var-requires */
import mousetrap from 'mousetrap'
import React, { useContext, useEffect, useRef } from 'react'

import { LocalConfigurationContext } from '../../context/LocalConfigurationContext'
import { Article } from '../models'
import styles from './ArticleContent.module.css'
import readable from './readable'

const getMql = () => window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)')

interface Props {
  article: Article
}

const getHTMLContent = (body: string, theme: string) => `
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <style>
      ${readable.css}
    </style>
  </head>
  <body data-theme="${theme}">
    ${body}
    <script>
window.onload = function() {
  document.querySelectorAll('a').forEach(a => a.setAttribute('target', '_blank'))
}
    </script>
  </body>
</html>
`

export default ({ article }: Props) => {
  const contentRef = useRef<HTMLDivElement>(null)
  const { localConfiguration } = useContext(LocalConfigurationContext)
  let { theme } = localConfiguration
  if (theme === 'auto') {
    const mql = getMql()
    theme = mql && mql.matches ? 'dark' : 'light'
  }

  useEffect(() => {
    if (contentRef.current) {
      const ifrm = document.createElement('iframe')
      contentRef.current.appendChild(ifrm)
      const doc = ifrm.contentWindow ? ifrm.contentWindow.document : ifrm.contentDocument
      if (doc) {
        doc.open()
        doc.write(getHTMLContent(article.html || article.text, theme))
        // Keyboard events have to propagate outside the iframe
        doc.onkeydown = function (e: KeyboardEvent) {
          switch (e.keyCode) {
            case 8:
              mousetrap.trigger('backspace')
              break
            case 77:
              mousetrap.trigger('m')
              break
            case 79:
              mousetrap.trigger('o')
              break
            case 82:
              mousetrap.trigger('r')
              break
            case 83:
              mousetrap.trigger('s')
              break
            case 191:
              mousetrap.trigger('?')
              break
            // default:
            // console.log(e.keyCode)
          }
        }
        doc.close()
      }
      ifrm.focus()
    }
  }, [article, theme])

  return <article className={styles.content} ref={contentRef} />
}
