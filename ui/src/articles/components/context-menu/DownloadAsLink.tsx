import React, { useCallback, useContext } from 'react'
import ReactModal from 'react-modal'
import { useModal } from 'react-modal-hook'

import LinkIcon from '../../../components/LinkIcon'
import styles from '../../../components/Dialog.module.css'
import { MessageContext } from '../../../context/MessageContext'
import fetchAPI from '../../../helpers/fetchAPI'
import { Article } from '../../models'
import DownloadPanel from './DownloadPanel'

interface Props {
  article: Article
  keyboard?: boolean
}

export default ({ article }: Props) => {
  const { showErrorMessage } = useContext(MessageContext)

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
  const [showDownloadModal, hideDownloadModal] = useModal(() => (
    <ReactModal
      isOpen
      shouldCloseOnEsc
      shouldCloseOnOverlayClick
      shouldFocusAfterRender
      appElement={document.getElementById('root')!}
      onRequestClose={hideDownloadModal}
      className={styles.dialog}
      overlayClassName={styles.overlay}
      style={{ content: { minWidth: '50vw' } }}
    >
      <DownloadPanel onCancel={hideDownloadModal} download={download} />
    </ReactModal>
  ))

  return (
    <LinkIcon title="Download article as ..." icon="download" onClick={showDownloadModal}>
      <span>Download as ...</span>
    </LinkIcon>
  )
}
