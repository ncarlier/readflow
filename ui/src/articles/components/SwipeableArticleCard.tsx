import React, { useCallback } from 'react'
import { useMutation } from 'react-apollo-hooks'

import { getGQLError } from '../../common/helpers'
import Icon from '../../common/Icon'
import SwipeableListItem from '../../common/SwipeableListItem'
import { connectMessageDispatch, IMessageDispatchProps } from '../../containers/MessageContainer'
import { Article, UpdateArticleStatusRequest } from '../models'
import { UpdateArticleStatus } from '../queries'
import ArticleCard from './ArticleCard'
import styles from './SwipeableArticleCard.module.css'

interface Props {
  article: Article
}

type AllProps = Props & IMessageDispatchProps

const Background = ({ icon }: { icon: string }) => (
  <div className={styles.background}>
    <Icon name={icon} />
  </div>
)

export const SwipeableArticleCard = (props: AllProps) => {
  const { article, showMessage } = props
  const updateArticleStatusMutation = useMutation<UpdateArticleStatusRequest>(UpdateArticleStatus)

  const updateArticleStatus = async (status: string) => {
    try {
      await updateArticleStatusMutation({
        variables: { id: article.id, status }
      })
    } catch (err) {
      showMessage(getGQLError(err), true)
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

export default connectMessageDispatch(SwipeableArticleCard)
