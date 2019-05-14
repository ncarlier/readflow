import React, { useCallback, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'

import ButtonIcon from '../../common/ButtonIcon'
import { getGQLError } from '../../common/helpers'
import { connectMessageDispatch, IMessageDispatchProps } from '../../containers/MessageContainer'
import useKeyboard from '../../hooks/useKeyboard'
import { Article, UpdateArticleStatusRequest } from '../models'
import { UpdateArticleStatus } from '../queries'

interface Props {
  article: Article
  floating?: boolean
  keyboard?: boolean
  onSuccess?: (article: Article) => void
}

type AllProps = Props & IMessageDispatchProps

export const MarkAsButton = (props: AllProps) => {
  const { article, floating = false, keyboard = false, showMessage, onSuccess } = props

  const [loading, setLoading] = useState(false)
  const updateArticleStatusMutation = useMutation<UpdateArticleStatusRequest>(UpdateArticleStatus)

  const updateArticleStatus = async (status: string) => {
    try {
      setLoading(true)
      await updateArticleStatusMutation({
        variables: { id: article.id, status }
        // update: updateCacheAfterUpdateStatus
      })
      if (floating) setLoading(false)
      if (onSuccess) onSuccess(article)
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
  useKeyboard('m', handleOnClick, keyboard)
  const kbs = keyboard ? ' [m]' : ''

  if (article.status === 'read') {
    return (
      <ButtonIcon
        title={'Mark as unread' + kbs}
        onClick={handleOnClick}
        loading={loading}
        floating={floating}
        icon="undo"
        primary
      />
    )
  }

  return (
    <ButtonIcon
      title={'Mark as read' + kbs}
      onClick={handleOnClick}
      loading={loading}
      floating={floating}
      icon="done"
      primary
    />
  )
}

export default connectMessageDispatch(MarkAsButton)
