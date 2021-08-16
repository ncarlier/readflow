import React, { useCallback, useContext } from 'react'
import { useMutation } from '@apollo/client'

import Icon from '../../components/Icon'
import SwipeableListItem from '../../components/SwipeableListItem'
import { MessageContext } from '../../contexts/MessageContext'
import { getGQLError } from '../../helpers'
import { updateCacheAfterUpdate } from '../cache'
import { Article, ArticleStatus, UpdateArticleRequest, UpdateArticleResponse } from '../models'
import { UpdateArticle } from '../queries'
import ArticleCard from './ArticleCard'
import styles from './SwipeableArticleCard.module.css'

interface Props {
  article: Article
}

const Background = ({ icon }: { icon: string }) => (
  <div className={styles.background}>
    <Icon name={icon} />
  </div>
)

export default (props: Props) => {
  const { article } = props

  const { showErrorMessage } = useContext(MessageContext)
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

  const handleOnDelete = useCallback(() => {
    const status = article.status === 'read' ? 'unread' : 'read'
    updateArticleStatus(status)
  }, [article, updateArticleStatus])

  const bgIcon = article.status === 'read' ? 'undo' : 'done'

  return (
    <SwipeableListItem background={<Background icon={bgIcon} />} onSwipe={handleOnDelete}>
      <ArticleCard article={article} isActive={false} />
    </SwipeableListItem>
  )
}
