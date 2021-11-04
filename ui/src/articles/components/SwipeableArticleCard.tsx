import React, { useCallback } from 'react'
import { useMutation } from '@apollo/client'

import { Icon, SwipeableListItem } from '../../components'
import { useMessage } from '../../contexts'
import { getGQLError } from '../../helpers'
import { updateCacheAfterUpdate } from '../cache'
import { Article, ArticleStatus, UpdateArticleRequest, UpdateArticleResponse } from '../models'
import { UpdateArticle } from '../queries'
import { ArticleCard } from '.'
import styles from './SwipeableArticleCard.module.css'

interface Props {
  article: Article
}

const Background = ({ icon }: { icon: string }) => (
  <div className={styles.background}>
    <Icon name={icon} />
  </div>
)

export const SwipeableArticleCard = (props: Props) => {
  const { article } = props

  const { showErrorMessage } = useMessage()
  const [updateArticleMutation] = useMutation<UpdateArticleResponse, UpdateArticleRequest>(UpdateArticle)

  const updateArticleStatus = useCallback(
    async (status: ArticleStatus) => {
      try {
        await updateArticleMutation({
          variables: { id: article.id, status },
          update: updateCacheAfterUpdate,
        })
      } catch (err) {
        showErrorMessage(getGQLError(err))
      }
    },
    [updateArticleMutation, article, showErrorMessage]
  )

  const handleOnSwipe = useCallback(() => {
    const status = article.status === 'read' ? 'inbox' : 'read'
    updateArticleStatus(status)
  }, [article, updateArticleStatus])

  const bgIcon = article.status === 'read' ? 'undo' : 'done'

  return (
    <SwipeableListItem background={<Background icon={bgIcon} />} onSwipe={handleOnSwipe}>
      <ArticleCard article={article} isActive={false} />
    </SwipeableListItem>
  )
}
