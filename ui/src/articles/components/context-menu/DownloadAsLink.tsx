import React, { useCallback, useContext, useState } from 'react'

import LinkIcon from '../../../components/LinkIcon'
import { MessageContext } from '../../../contexts/MessageContext'
import fetchAPI from '../../../helpers/fetchAPI'
import { Article } from '../../models'
import DownloadPanel from './DownloadPanel'
import Overlay from '../../../components/Overlay'
import DownloadProgress from './DownloadProgress'

interface Props {
  article: Article
  keyboard?: boolean
}

export default ({ article }: Props) => {
  const { showErrorMessage } = useContext(MessageContext)
  const [isVisible, setIsVisible] = useState(false)
  const [isDownloading, setIsDownloading] = useState(false)
  const [contentLength, setContentLength] = useState(0)
  const [contentReceived, setContentReceived] = useState(0)
  const showOverlay = () => setIsVisible(true)
  const hideOverlay = () => setIsVisible(false)

  const download = useCallback(
    async (format: string) => {
      try {
        const res = await fetchAPI(`/articles/${article.id}`, { f: format }, { method: 'GET' })
        if (res.ok && res.body) {
          const reader = res.body.getReader()
          const contentLength = parseInt(res.headers.get('X-Content-Length') || '0')
          setContentLength(contentLength)
          if (contentLength > 1_048_576) {
            setIsDownloading(true)
          }
          const contentType = res.headers.get('Content-Type') || 'application/octet-stream'
          let receivedLength = 0
          const chunks = []
          for (;;) {
            const { done, value } = await reader.read()
            if (done) {
              break
            }
            if (value) {
              chunks.push(value)
              receivedLength += value.length
            }
            setContentReceived(receivedLength)
          }
          const chunksAll = new Uint8Array(receivedLength)
          let position = 0
          for (const chunk of chunks) {
            chunksAll.set(chunk, position)
            position += chunk.length
          }
          const data = new Blob([chunksAll.buffer], { type: contentType })
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
      } finally {
        setIsDownloading(false)
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
      <Overlay visible={isDownloading}>
        <DownloadProgress total={contentLength} value={contentReceived} />
      </Overlay>
    </>
  )
}
