import React, { useRef, useEffect } from 'react'

import { Article } from '../models'

import styles from './ArticleContent.module.css'

type Props = {
  article: Article
}

export default ({article}: Props) => {
  const contentRef = useRef<any>(null)

  var cssLink = document.createElement("link")
  cssLink.href = process.env.PUBLIC_URL + "/readable.css"
  cssLink.rel = "stylesheet" 
  cssLink.type = "text/css" 
 
  useEffect(
    () => {
      let ifrm = contentRef.current
      ifrm = ifrm.contentWindow || ifrm.contentDocument.document || ifrm.contentDocument
      ifrm.document.open()
      ifrm.document.write(article.html)
      ifrm.document.head.appendChild(cssLink)
      ifrm.document.close()
    },
    [article],
  )

  return (
    <article className={styles.content}>
      <iframe ref={contentRef}></iframe>
    </article>
  )
}
