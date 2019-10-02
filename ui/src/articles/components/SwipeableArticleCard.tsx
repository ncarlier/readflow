import React, { useCallback, useContext } from 'react'
import { useMutation } from 'react-apollo-hooks'

import { getGQLError } from '../../helpers'
import Icon from '../../components/Icon'
import SwipeableListItem from '../../components/SwipeableListItem'
import { MessageContext } from '../../context/MessageContext'
import { Article, UpdateArticleStatusRequest } from '../models'
import { UpdateArticleStatus } from '../queries'
import styles from './SwipeableArticleCard.module.css'
import ArticleCard from './ArticleCard'

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
  const updateArticleStatusMutation = useMutation<UpdateArticleStatusRequest>(UpdateArticleStatus)

  const updateArticleStatus = async (status: string) => {
    try {
      await updateArticleStatusMutation({
        variables: { id: article.id, status }
      })
    } catch (err) {
      showErrorMessage(getGQLError(err))
    }
  }

  const handleOnDelete = useCallback(() => {
    const status = article.status === 'read' ? 'unread' : 'read'
    updateArticleStatus(status)
  }, [article])

  const bgIcon = article.status === 'read' ? 'undo' : 'done'

  return (
    <SwipeableListItem background={<Background icon={bgIcon} />} onSwipe={handleOnDelete}>
      <ArticleCard article={article} isActive={false} />
    </SwipeableListItem>
  )
}
