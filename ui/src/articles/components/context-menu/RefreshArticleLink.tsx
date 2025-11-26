import React, { SyntheticEvent, useCallback } from 'react'
import { useMutation } from '@apollo/client'

import { LinkIcon } from '../../../components/LinkIcon'
import { Kbd } from '../../../components'
import { useMessage } from '../../../contexts'
import { getGQLError } from '../../../helpers'
import { Article, UpdateArticleRequest, UpdateArticleResponse } from '../../models'
import { UpdateFullArticle } from '../../queries'
import { updateCacheAfterUpdate } from '../../cache'

interface Props {
  article: Article
  keyboard?: boolean
  onLoading: (loading: boolean) => void
}

export const RefreshArticleLink = ({ article, keyboard, onLoading }: Props) => {
  const { showMessage, showErrorMessage } = useMessage()
  const [updateArticleMutation] = useMutation<UpdateArticleResponse, UpdateArticleRequest>(UpdateFullArticle)

  const refreshArticle = useCallback(
    async () => {
      onLoading(true)
      try {
        const variables = { id: article.id, refresh: true }
        const res = await updateArticleMutation({
          variables,
          update: updateCacheAfterUpdate,
        })
        if (res.data) {
          const { article: updatedArticle } = res.data.updateArticle
          showMessage(`Article "${updatedArticle.title}" (#${updatedArticle.id}) refreshed`)
        }
      } catch (err) {
        showErrorMessage(getGQLError(err))
      } finally {
        onLoading(false)
      }
    },
    [article, updateArticleMutation, showMessage, showErrorMessage]
  )
  
  const handleOnClick = useCallback(
    (ev: SyntheticEvent) => {
      ev.preventDefault()
      refreshArticle()
    },
    [refreshArticle]
  )
  
  return (
    <LinkIcon title="Refresh article content" icon="refresh" onClick={handleOnClick}>
      <span>Refresh content</span>
      {keyboard && <Kbd keys="alt+r" onKeypress={refreshArticle} />}
    </LinkIcon>
  )
}
