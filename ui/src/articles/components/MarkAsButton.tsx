
import React, { useState, useCallback } from 'react'
import { useMutation } from 'react-apollo-hooks'

import { Article, UpdateArticleStatusRequest } from '../models'
import ButtonIcon from '../../common/ButtonIcon'

import { UpdateArticleStatus } from '../queries'
import { getGQLError } from '../../common/helpers'
import { connectMessageDispatch, IMessageDispatchProps } from '../../containers/MessageContainer';
import useKeyboard from '../../hooks/useKeyboard';

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
  const updateArticleStatusMutation = useMutation<UpdateArticleStatusRequest>(UpdateArticleStatus)
  
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

  // Keyboard shortcut is only active for Floating Action Button
  useKeyboard('m', handleOnClick, floating)
  const kbs = floating ? " [m]" : ""

  if (article.status === 'read') {
    return (
      <ButtonIcon
        title={"Mark as unread" + kbs}
        onClick={handleOnClick}
        loading={loading}
        floating={floating}
        icon="undo"
        primary />
    )
  }

  return (
    <ButtonIcon
      title={"Mark as read" + kbs}
      onClick={handleOnClick}
      loading={loading}
      floating={floating}
      icon="done"
      primary />
  )
}

export default connectMessageDispatch(MarkAsButton)
