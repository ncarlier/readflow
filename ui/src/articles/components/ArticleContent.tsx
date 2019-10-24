/* eslint-disable @typescript-eslint/no-var-requires */
import mousetrap from 'mousetrap'
import React, { useEffect, useRef } from 'react'

import { Article } from '../models'
import styles from './ArticleContent.module.css'

const css = require('!!raw-loader!./readable.css')

interface Props {
  article: Article
}

const getHTMLContent = (body: string) => `
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <style>
      ${css.default}
    </style>
  </head>
  <body>
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

  useEffect(() => {
    if (contentRef.current) {
      var ifrm = document.createElement('iframe')
      contentRef.current.appendChild(ifrm)
      let doc = ifrm.contentWindow ? ifrm.contentWindow.document : ifrm.contentDocument
      if (doc) {
        doc.open()
        doc.write(getHTMLContent(article.html || article.text))
        // Keyboard events have to propagate outside the iframe
        doc.onkeydown = function(e: KeyboardEvent) {
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
  }, [])

  return <article className={styles.content} ref={contentRef} />
}
