import React, { useCallback, useContext, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'

import ButtonIcon from '../../components/ButtonIcon'
import { MessageContext } from '../../context/MessageContext'
import { getGQLError } from '../../helpers'
import useKeyboard from '../../hooks/useKeyboard'
import { updateCacheAfterUpdate } from '../cache'
import { Article, UpdateArticleRequest, UpdateArticleResponse } from '../models'
import { UpdateArticle } from '../queries'

interface Props {
  article: Article
  keyboard?: boolean
  onSuccess?: (article: Article) => void
}

export default (props: Props) => {
  const { article, keyboard = false, onSuccess } = props

  const { showErrorMessage } = useContext(MessageContext)
  const [loading, setLoading] = useState(false)
  const [updateArticleMutation] = useMutation<UpdateArticleResponse, UpdateArticleRequest>(UpdateArticle)

  const updateArticle = useCallback(
    async (starred: boolean) => {
      try {
        setLoading(true)
        await updateArticleMutation({
          variables: { id: article.id, starred },
          update: updateCacheAfterUpdate
        })
        setLoading(false)
        if (onSuccess) onSuccess(article)
      } catch (err) {
        setLoading(false)
        showErrorMessage(getGQLError(err))
      }
    },
    [updateArticleMutation, article, onSuccess, showErrorMessage]
  )

  const handleOnClick = useCallback(() => {
    updateArticle(!article.starred)
  }, [article, updateArticle])

  // Keyboard shortcut is only active for Floating Action Button
  useKeyboard('s', handleOnClick, keyboard)
  const kbs = keyboard ? ' [s]' : ''

  return (
    <ButtonIcon
      title={article.starred ? `Unstars this article${kbs}` : `Stars this article${kbs}`}
      style={article.starred ? { color: 'gold' } : undefined}
      onClick={handleOnClick}
      loading={loading}
      icon={article.starred ? 'star' : 'star_outline'}
    />
  )
}
