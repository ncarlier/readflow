import React, { useCallback, useContext, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'

import ButtonIcon from '../../components/ButtonIcon'
import { getGQLError } from '../../helpers'
import { MessageContext } from '../../context/MessageContext'
import useKeyboard from '../../hooks/useKeyboard'
import { Article, UpdateArticleStatusRequest } from '../models'
import { UpdateArticleStatus } from '../queries'

interface Props {
  article: Article
  floating?: boolean
  keyboard?: boolean
  onSuccess?: (article: Article) => void
}

export default (props: Props) => {
  const { article, floating = false, keyboard = false, onSuccess } = props

  const { showErrorMessage } = useContext(MessageContext)
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
      showErrorMessage(getGQLError(err))
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
        variant="primary"
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
      variant="primary"
    />
  )
}
