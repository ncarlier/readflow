
import React from 'react'
import { useApolloClient } from 'react-apollo-hooks'

import {Article, GetArticleResponse} from '../models'

import { GetFullArticle } from '../queries'
import { connectMessageDispatch, IMessageDispatchProps } from '../../containers/MessageContainer'
import { connectOfflineDispatch, IOfflineDispatchProps } from '../../containers/OfflineContainer'
import ConfirmDialog from '../../common/ConfirmDialog'
import { useModal } from 'react-modal-hook'
import LinkIcon from '../../common/LinkIcon'
import useKeyboard from '../../hooks/useKeyboard'

type Props = {
  article: Article
  remove?: boolean
  noShortcuts?: boolean
}

type AllProps = Props & IMessageDispatchProps & IOfflineDispatchProps

export const OfflineLink = (props: AllProps) => {
  const {
    article,
    remove,
    noShortcuts,
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
  
  useKeyboard('r', showDeleteConfirmModal, !noShortcuts && remove)
  useKeyboard('o', putArticleOffline, !noShortcuts && !remove)

  if (remove) {
    return (
      <LinkIcon
        title="Remove"
        onClick={showDeleteConfirmModal}
        icon="delete">
        <span>Remove offline</span>{!noShortcuts && <kbd>r</kbd>}
      </LinkIcon>
    )
  }

  return (
    <LinkIcon
      title="Put offline"
      onClick={putArticleOffline}
      icon="signal_wifi_off">
      <span>Put offline</span>{!noShortcuts && <kbd>o</kbd>}
    </LinkIcon>
  )
}

export default connectOfflineDispatch(
  connectMessageDispatch(OfflineLink)
)
