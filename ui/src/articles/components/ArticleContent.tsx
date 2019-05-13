import mousetrap from 'mousetrap'
import React, { useEffect, useRef } from 'react'

import { Article } from '../models'
import styles from './ArticleContent.module.css'

interface Props {
  article: Article
}

export default ({ article }: Props) => {
  const contentRef = useRef<HTMLDivElement>(null)

  var cssLink = document.createElement('link')
  cssLink.href = process.env.PUBLIC_URL + '/readable.css'
  cssLink.rel = 'stylesheet'
  cssLink.type = 'text/css'
  var script = document.createElement('script')
  script.setAttribute('type', 'text/javascript')
  script.setAttribute('src', process.env.PUBLIC_URL + '/readable.js')

  useEffect(() => {
    if (contentRef.current) {
      var ifrm = document.createElement('iframe')
      contentRef.current.appendChild(ifrm)
      let doc = ifrm.contentWindow ? ifrm.contentWindow.document : ifrm.contentDocument
      if (doc) {
        doc.open()
        doc.write(article.html)
        doc.head.appendChild(cssLink)
        doc.head.appendChild(script)
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
