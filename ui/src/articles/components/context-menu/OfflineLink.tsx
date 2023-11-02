import React, { useCallback, useState } from 'react'
import { useModal } from 'react-modal-hook'

import { ConfirmDialog, Kbd, LinkIcon, Loader, Overlay } from '../../../components'
import { connectOffline, OfflineProps } from '../../../containers/OfflineContainer'
import { useMessage } from '../../../contexts'
import { Article } from '../../models'
import { useAPI } from '../../../hooks'
import { fetchAPI } from '../../../helpers'

interface Props {
  article: Article
  keyboard?: boolean
}

type AllProps = Props & OfflineProps

const blobToDataURL = async (blob: Blob): Promise<string> => new Promise((resolve) => {
  const reader = new FileReader()
  reader.onloadend = () => resolve(reader.result as string)
  reader.readAsDataURL(blob)
})

export const OfflineLink = (props: AllProps) => {
  const { article, keyboard = false, saveOfflineArticle, removeOfflineArticle, offlineArticles } = props
  const fetchArticleContent = useAPI(`/articles/${article.id}`, { method: 'GET' })
  const { showMessage, showErrorMessage } = useMessage()
  const [loading, setLoading] = useState(false)

  const putArticleOffline = useCallback(async () => {
    setLoading(true)
    const offlineArticle = { ...article }
    try {
      // download article content with embedded images
      const res = await fetchArticleContent({ f: 'html-single' })
      if (res) {
        if (res.ok && res.body) {
          offlineArticle.html = await res.text()
          if (article.image && article.thumbnails) {
            const thumbnail = [...article.thumbnails].sort((a, b) => parseInt(b.size) - parseInt(a.size))[0]
            // download article image
            const img = await fetchAPI(`/img/${thumbnail.hash}/resize:fit:${thumbnail.size}/${btoa(article.image)}`)
            if (img && img.ok && img.body) {
              const blob = await img.blob()
              offlineArticle.image = await blobToDataURL(blob)
            }
          }
          await saveOfflineArticle(offlineArticle)
          showMessage(`Article put offline: ${article.title}`)
        } else {
          const err = await res.json()
          throw new Error(err.detail || res.statusText)
        }
      }
    } catch (err: any) {
      showErrorMessage(err.message)
    } finally {
      setLoading(false)
    }
  }, [article, fetchArticleContent, saveOfflineArticle, showMessage])

  const deleteArticleOffline = useCallback(async () => {
    try {
      await removeOfflineArticle(article)
      showMessage(`Article deleted from offline storage: ${article.title}`)
    } catch (err: any) {
      showErrorMessage(err.message)
    }
  }, [article, removeOfflineArticle, showMessage, showErrorMessage])

  const [showDeleteConfirmModal, hideDeleteConfirmModal] = useModal(() => (
    <ConfirmDialog
      title={article.title}
      confirmLabel="Delete"
      onConfirm={() => deleteArticleOffline()}
      onCancel={hideDeleteConfirmModal}
    >
      Deleting an article from offline storage is irreversible. Please confirm!
    </ConfirmDialog>
  ))

  if (article.isOffline) {
    return (
      <LinkIcon title="Delete" onClick={showDeleteConfirmModal} icon="delete">
        <span>Delete offline</span>
        {keyboard && <Kbd keys="d" onKeypress={showDeleteConfirmModal} />}
      </LinkIcon>
    )
  }

  return (
    <>
      <LinkIcon title="Put offline" onClick={putArticleOffline} icon="signal_wifi_off">
        <span>Put offline</span>
        {keyboard && <Kbd keys="o" onKeypress={putArticleOffline} />}
      </LinkIcon>
      <Overlay visible={loading || offlineArticles.loading}>
        <Loader center />
      </Overlay>
    </>
  )
}

export default connectOffline(OfflineLink)
