import React, { useCallback, useContext, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'

import ButtonIcon from '../../components/ButtonIcon'
import { MessageContext } from '../../context/MessageContext'
import { getGQLError } from '../../helpers'
import useKeyboard from '../../hooks/useKeyboard'
import { updateCacheAfterUpdate } from '../cache'
import { Article, ArticleStatus, UpdateArticleRequest, UpdateArticleResponse } from '../models'
import { UpdateArticle } from '../queries'

interface Props {
  article: Article
  floating?: boolean
  keyboard?: boolean
  onSuccess?: (article: Article) => void
}

export default (props: Props) => {
  const { article, floating = false, keyboard = false, onSuccess } = props

  const { showErrorMessage } = useContext(MessageContext)
  const [loading, setLoading] = useState(false)
  const [updateArticleMutation] = useMutation<UpdateArticleResponse, UpdateArticleRequest>(UpdateArticle)

  const updateArticleStatus = useCallback(
    async (status: ArticleStatus) => {
      try {
        setLoading(true)
        await updateArticleMutation({
          variables: { id: article.id, status },
          update: updateCacheAfterUpdate
        })
        if (!floating) setLoading(false)
        if (onSuccess) onSuccess(article)
      } catch (err) {
        setLoading(false)
        showErrorMessage(getGQLError(err))
      }
    },
    [updateArticleMutation, article, floating, onSuccess, showErrorMessage]
  )

  const handleOnClick = useCallback(() => {
    const status = article.status === 'read' ? 'unread' : 'read'
    updateArticleStatus(status)
  }, [article, updateArticleStatus])

  // Keyboard shortcut is only active for Floating Action Button
  useKeyboard('m', handleOnClick, keyboard)
  const kbs = keyboard ? ' [m]' : ''

  if (article.status === 'read') {
    return (
      <ButtonIcon
        title={'Mark as unread' + kbs}
        onClick={handleOnClick}
        loading={loading}
        floating={floating}
        icon="undo"
        variant="primary"
      />
    )
  }

  return (
    <ButtonIcon
      title={'Mark as read' + kbs}
      onClick={handleOnClick}
      loading={loading}
      floating={floating}
      icon="done"
      variant="primary"
    />
  )
}
