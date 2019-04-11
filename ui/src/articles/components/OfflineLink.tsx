
import React, { useState, ReactNode } from 'react'
import { useApolloClient } from 'react-apollo-hooks'

import {Article, GetArticleResponse} from '../models'
import ButtonIcon from '../../common/ButtonIcon'

import { GetFullArticle } from '../queries'
import { connectMessageDispatch, IMessageDispatchProps } from '../../containers/MessageContainer'
import { connectOfflineDispatch, IOfflineDispatchProps } from '../../containers/OfflineContainer'
import ConfirmDialog from '../../common/ConfirmDialog'
import { useModal } from 'react-modal-hook'
import LinkIcon from '../../common/LinkIcon'

type Props = {
  article: Article
  remove?: boolean
}

type AllProps = Props & IMessageDispatchProps & IOfflineDispatchProps

export const OfflineLink = (props: AllProps) => {
  const {
    article,
    remove,
    saveOfflineArticle,
    removeOfflineArticle,
    showMessage
  } = props

  const client = useApolloClient()
  
  const putArticleOffline = async () => {
    try {
      const { errors, data } = await client.query<GetArticleResponse>({
        query: GetFullArticle,
        variables: {id: article.id}
      })
      if (data) {
        const fullArticle = { ...article, ...data.article} 
        await saveOfflineArticle(fullArticle)
        showMessage(`Article put offline: ${article.title}`)
      }
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
      <LinkIcon
        title="Remove"
        onClick={showDeleteConfirmModal}
        icon="delete">
        <span>Remove offline</span><small>[r]</small>
      </LinkIcon>
    )
  }

  return (
    <LinkIcon
      title="Put offline"
      onClick={putArticleOffline}
      icon="signal_wifi_off">
      <span>Put offline</span><small>[o]</small>
    </LinkIcon>
  )
}

export default connectOfflineDispatch(
  connectMessageDispatch(OfflineLink)
)
