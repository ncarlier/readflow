import React, { useCallback, useContext, useEffect, useState } from 'react'
import { useMutation } from '@apollo/client'

import { ButtonIcon } from '../../components'
import { MessageContext } from '../../contexts/MessageContext'
import { getGQLError } from '../../helpers'
import { useKeyboard } from '../../hooks'
import { updateCacheAfterUpdate } from '../cache'
import { Article, ArticleStatus, UpdateArticleRequest, UpdateArticleResponse } from '../models'
import { UpdateArticle } from '../queries'

interface Props {
  article: Article
  floating?: boolean
  keyboard?: boolean
  onSuccess?: (article: Article) => void
}

export const MarkAsButton = (props: Props) => {
  const isMounted = React.useRef(true)
  const { article, floating = false, keyboard = false, onSuccess } = props

  const { showErrorMessage } = useContext(MessageContext)
  const [loading, setLoading] = useState(false)
  const [updateArticleMutation] = useMutation<UpdateArticleResponse, UpdateArticleRequest>(UpdateArticle)

  // Small tips to prevent update warnings on unmounted components
  useEffect(
    () => () => {
      isMounted.current = false
    },
    []
  )

  const updateArticleStatus = useCallback(
    async (status: ArticleStatus) => {
      setLoading(true)
      try {
        await updateArticleMutation({
          variables: { id: article.id, status },
          update: updateCacheAfterUpdate,
        })
        if (onSuccess) onSuccess(article)
      } catch (err) {
        showErrorMessage(getGQLError(err))
      } finally {
        if (isMounted.current) {
          setLoading(false)
        }
      }
    },
    [updateArticleMutation, article, onSuccess, showErrorMessage]
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
