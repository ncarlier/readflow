import React, { useCallback, useContext, useState } from 'react'

import LinkIcon from '../../../components/LinkIcon'
import { MessageContext } from '../../../contexts/MessageContext'
import fetchAPI from '../../../helpers/fetchAPI'
import { Article } from '../../models'
import DownloadPanel from './DownloadPanel'
import Overlay from '../../../components/Overlay'

interface Props {
  article: Article
  keyboard?: boolean
}

export default ({ article }: Props) => {
  const { showErrorMessage } = useContext(MessageContext)
  const [isVisible, setIsVisible] = useState(false)
  const showOverlay = () => setIsVisible(true)
  const hideOverlay = () => setIsVisible(false)

  const download = useCallback(
    async (format: string) => {
      try {
        const res = await fetchAPI(`/articles/${article.id}`, { f: format }, { method: 'GET' })
        if (res.ok) {
          const data = await res.blob()
          const href = window.URL.createObjectURL(data)
          const link = document.createElement('a')
          link.href = href
          link.setAttribute('download', `${article.title}.html`)
          document.body.appendChild(link)
          link.click()
          document.body.removeChild(link)
        } else {
          const err = await res.json()
          throw new Error(err.detail || res.statusText)
        }
      } catch (err) {
        showErrorMessage(err.message)
      }
    },
    [article, showErrorMessage]
  )

  return (
    <>
      <LinkIcon title="Download article as ..." icon="download" onClick={showOverlay}>
        <span>Download as ...</span>
      </LinkIcon>
      <Overlay visible={isVisible}>
        <DownloadPanel onCancel={hideOverlay} download={download} />
      </Overlay>
    </>
  )
}
