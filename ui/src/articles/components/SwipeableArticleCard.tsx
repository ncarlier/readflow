import React, { useCallback } from 'react'
import SwipeToDelete from 'react-swipe-to-delete-component'
import 'react-swipe-to-delete-component/dist/swipe-to-delete.css'

import styles from './SwipeableArticleCard.module.css'
import {Article, UpdateArticleStatusRequest} from '../models'
import ArticleCard from './ArticleCard'
import Icon from '../../common/Icon'
import { useMutation } from 'react-apollo-hooks';
import { UpdateArticleStatus } from '../queries';
import { getGQLError } from '../../common/helpers';
import { IMessageDispatchProps, connectMessageDispatch } from '../../containers/MessageContainer';

type Props = {
  article: Article
  readMoreBasePath: string
}

type AllProps = Props & IMessageDispatchProps

const Background = ({icon}: {icon: string}) => (
  <div className={styles.background}>
    <Icon name={icon} />
    <span />
    <Icon name={icon} />
  </div>
)

export const SwipeableArticleCard = (props: AllProps) => {
  const {article, readMoreBasePath, showMessage} = props
  const updateArticleStatusMutation = useMutation<UpdateArticleStatusRequest>(UpdateArticleStatus)
  
  const updateArticleStatus = async (status: string) => {
    try{
      const res = await updateArticleStatusMutation({
        variables: {id: article.id, status},
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
    <SwipeToDelete
      background={<Background icon={bgIcon} />}
      deleteSwipe={0.3}
      onDelete={handleOnDelete}>
      <ArticleCard article={article} readMoreBasePath={readMoreBasePath} />
    </SwipeToDelete>
  )
}

export default connectMessageDispatch(SwipeableArticleCard)
