import React, { useCallback, useContext } from 'react'
import { useMutation } from 'react-apollo-hooks'

import Icon from '../../components/Icon'
import SwipeableListItem from '../../components/SwipeableListItem'
import { MessageContext } from '../../context/MessageContext'
import { getGQLError } from '../../helpers'
import { Article, UpdateArticleRequest } from '../models'
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
  const [updateArticleMutation] = useMutation<UpdateArticleRequest>(UpdateArticle)

  const updateArticleStatus = useCallback(
    async (status: string) => {
      try {
        await updateArticleMutation({
          variables: { id: article.id, status }
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
