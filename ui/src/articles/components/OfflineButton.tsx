
import React, { useState } from 'react'
import { useApolloClient } from 'react-apollo-hooks'

import {Article, GetArticleResponse} from '../models'
import ButtonIcon from '../../common/ButtonIcon'

import { GetFullArticle } from '../queries'
import { connectMessageDispatch, IMessageDispatchProps } from '../../containers/MessageContainer'
import { connectOfflineDispatch, IOfflineDispatchProps } from '../../containers/OfflineContainer'
import ConfirmDialog from '../../common/ConfirmDialog'
import { useModal } from 'react-modal-hook'

type Props = {
  article: Article
  remove?: boolean
}

type AllProps = Props & IMessageDispatchProps & IOfflineDispatchProps

export const OfflineButton = (props: AllProps) => {
  const {
    article,
    remove,
    saveOfflineArticle,
    removeOfflineArticle,
    showMessage
  } = props

  const [loading, setLoading] = useState(false)
  const client = useApolloClient()
  
  const putArticleOffline = async () => {
    try {
      setLoading(true)
      const { errors, data } = await client.query<GetArticleResponse>({
        query: GetFullArticle,
        variables: {id: article.id}
      })
      if (data) {
        const fullArticle = { ...article, ...data.article} 
        await saveOfflineArticle(fullArticle)
        showMessage(`Article put offline: ${article.title}`)
      }
      setLoading(false)
      if (errors) {
        throw new Error(errors[0])
      }
    } catch (err) {
      showMessage(err.message, true)
    }
  }

  const removeArticleOffline = async () => {
    try {
      const res = await removeOfflineArticle(article)
      showMessage(`Article removed from offline storage: ${article.title}`)
    } catch (err) {
      showMessage(err.message, true)
    } 
  }

  const [showDeleteConfirmModal, hideDeleteConfirmModal] = useModal(
    () => (
      <ConfirmDialog
        title="Remove article?"
        confirmLabel="Remove"
        onConfirm={() => removeArticleOffline()}
        onCancel={hideDeleteConfirmModal}
      >
        Removing an article from offline storage is irreversible. Please confirm!
      </ConfirmDialog>
    )
  )

  if (remove) {
    return (
      <ButtonIcon
        title="Remove"
        onClick={showDeleteConfirmModal}
        icon="delete"
        loading={loading} />
    )
  }

  return (
    <ButtonIcon
      title="Put offline"
      onClick={putArticleOffline}
      loading={loading}
      icon="signal_wifi_off"
    />
  )
}

export default connectOfflineDispatch(
  connectMessageDispatch(OfflineButton)
)
