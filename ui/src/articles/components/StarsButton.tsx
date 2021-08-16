import React, { useCallback, useContext, useState } from 'react'
import { useMutation } from '@apollo/client'

import { MessageContext } from '../../contexts/MessageContext'
import { getGQLError } from '../../helpers'
import useKeyboard from '../../hooks/useKeyboard'
import { updateCacheAfterUpdate } from '../cache'
import { Article, UpdateArticleRequest, UpdateArticleResponse } from '../models'
import { UpdateArticle } from '../queries'
import DrawerMenu from '../../components/DrawerMenu'
import Stars from '../../components/Stars'

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
    async (stars: number) => {
      try {
        setLoading(true)
        await updateArticleMutation({
          variables: { id: article.id, stars },
          update: updateCacheAfterUpdate,
        })
        if (onSuccess) onSuccess(article)
      } catch (err) {
        showErrorMessage(getGQLError(err))
      } finally {
        setLoading(false)
      }
    },
    [updateArticleMutation, article, onSuccess, showErrorMessage]
  )

  const handleOnKeyboard = useCallback(() => {
    updateArticle(article.stars > 0 ? 0 : 1)
  }, [article, updateArticle])

  // Keyboard shortcut is only active for Floating Action Button
  useKeyboard('s', handleOnKeyboard, keyboard)
  const kbs = keyboard ? ' [s]' : ''
  const title = `Star this article${kbs}`
  const style = article.stars > 0 ? { color: 'gold' } : undefined
  const icon = article.stars > 0 ? 'star' : 'star_outline'

  return (
    <DrawerMenu title={title} icon={loading ? 'loop' : icon} style={style}>
      <div style={{ textAlign: 'center' }}>
        <Stars value={article.stars} onChange={updateArticle} />
      </div>
    </DrawerMenu>
  )
}
