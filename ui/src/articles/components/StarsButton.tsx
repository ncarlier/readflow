import React, { useCallback, useContext, useState } from 'react'
import { useMutation } from '@apollo/client'

import { MessageContext } from '../../context/MessageContext'
import { getGQLError } from '../../helpers'
import useKeyboard from '../../hooks/useKeyboard'
import { updateCacheAfterUpdate } from '../cache'
import { Article, UpdateArticleRequest, UpdateArticleResponse } from '../models'
import { UpdateArticle } from '../queries'
import DropdownMenu, { DropDownOrigin } from '../../components/DropdownMenu'
import Stars from '../../components/Stars'

interface Props {
  article: Article
  keyboard?: boolean
  origin?: DropDownOrigin
  onSuccess?: (article: Article) => void
}

export default (props: Props) => {
  const { article, keyboard = false, origin, onSuccess } = props

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
    <DropdownMenu title={title} origin={origin} icon={loading ? 'loop' : icon} style={style}>
      <ul>
        <li>
          <Stars value={article.stars} onChange={updateArticle} />
        </li>
      </ul>
    </DropdownMenu>
  )
}
