import React from 'react'
import { useApolloClient } from '@apollo/client'
import { useModal } from 'react-modal-hook'

import { ConfirmDialog, Kbd, LinkIcon, Loader, Overlay } from '../../../components'
import { connectOffline, OfflineProps } from '../../../containers/OfflineContainer'
import { useMessage } from '../../../contexts'
import { Article, GetArticleResponse } from '../../models'
import { GetFullArticle } from '../../queries'

interface Props {
  article: Article
  keyboard?: boolean
}

type AllProps = Props & OfflineProps

export const OfflineLink = (props: AllProps) => {
  const { article, keyboard = false, saveOfflineArticle, removeOfflineArticle, offlineArticles } = props
  const { showMessage, showErrorMessage } = useMessage()
  const { loading } = offlineArticles

  const client = useApolloClient()

  const putArticleOffline = async () => {
    try {
      const { errors, data } = await client.query<GetArticleResponse>({
        query: GetFullArticle,
        variables: { id: article.id },
      })
      if (data) {
        const fullArticle = { ...article, ...data.article }
        await saveOfflineArticle(fullArticle)
        showMessage(`Article put offline: ${article.title}`)
      }
      if (errors) {
        throw new Error(errors[0].message)
      }
    } catch (err) {
      showErrorMessage(err.message)
    }
  }

  const deleteArticleOffline = async () => {
    try {
      await removeOfflineArticle(article)
      showMessage(`Article deleted from offline storage: ${article.title}`)
    } catch (err) {
      showErrorMessage(err.message)
    }
  }

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
      <Overlay visible={loading}>
        <Loader center />
      </Overlay>
    </>
  )
}

export default connectOffline(OfflineLink)
