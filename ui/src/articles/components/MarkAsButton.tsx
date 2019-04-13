
import React, { useState, useCallback } from 'react'
import { useMutation } from 'react-apollo-hooks'

import { Article } from '../models'
import ButtonIcon from '../../common/ButtonIcon'

import { UpdateArticleStatus } from '../queries'
import { getGQLError } from '../../common/helpers'
import { updateCacheAfterUpdateStatus } from '../cache'
import { connectMessageDispatch, IMessageDispatchProps } from '../../containers/MessageContainer';

type UpdateArticleStatusFields = {
  id: number
  status: string
}

type Props = {
  article: Article
  floating?: boolean
}

type AllProps = Props & IMessageDispatchProps

export const MarkAsButton = (props: AllProps) => {
  const {
    article,
    floating = false,
    showMessage
  } = props

  const [loading, setLoading] = useState(false)
  const updateArticleStatusMutation = useMutation<UpdateArticleStatusFields>(UpdateArticleStatus)
  
  const updateArticleStatus = async (status: string) => {
    try{
      setLoading(true)
      const res = await updateArticleStatusMutation({
        variables: {id: article.id, status},
        // update: updateCacheAfterUpdateStatus
      })
      if (floating) setLoading(false)
    } catch (err) {
      setLoading(false)
      showMessage(getGQLError(err), true)
    }
  }

  const handleOnClick = useCallback(() => {
    const status = article.status === 'read' ? 'unread' : 'read'
    updateArticleStatus(status) 
  }, [article])

  if (article.status === 'read') {
    return (
      <ButtonIcon
        title="Mark as unread"
        onClick={handleOnClick}
        loading={loading}
        floating={floating}
        icon="undo"
        primary />
    )
  }

  return (
    <ButtonIcon
      title="Mark as read"
      onClick={handleOnClick}
      loading={loading}
      floating={floating}
      icon="done"
      primary />
  )
}

export default connectMessageDispatch(MarkAsButton)
