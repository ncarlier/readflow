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
        doc.close()
      }
    }
  }, [])

  return <article className={styles.content} ref={contentRef} />
}
