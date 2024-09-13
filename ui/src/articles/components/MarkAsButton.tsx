import React, { useCallback, useEffect, useState } from 'react'
import { useMutation } from '@apollo/client'

import { ButtonIcon } from '../../components'
import { useMessage } from '../../contexts'
import { getGQLError } from '../../helpers'
import { useKeyboard } from '../../hooks'
import { updateCacheAfterUpdate } from '../cache'
import { Article, ArticleStatus, UpdateArticleRequest, UpdateArticleResponse } from '../models'
import { UpdateArticle } from '../queries'

interface Props {
  article: Article
  status: ArticleStatus
  floating?: boolean
  keyboard?: boolean
  onSuccess?: (article: Article) => void
}

type VariantsType = {
  [K in ArticleStatus]: {
    title: string
    icon: string
    kbs: string
  }
}

const variants: VariantsType = {
  inbox: {
    title: 'Put back to the inbox',
    icon: 'undo',
    kbs: 'ins',
  },
  read: {
    title: 'Mark as read',
    icon: 'done',
    kbs: 'del',
  },
  to_read: {
    title: 'Read it later',
    icon: 'book',
    kbs: 'r',
  },
}

export const MarkAsButton = (props: Props) => {
  const isMounted = React.useRef(true)
  const { article, status, floating = false, keyboard = false, onSuccess } = props

  const { showErrorMessage } = useMessage()
  const [loading, setLoading] = useState(false)
  const [updateArticleMutation] = useMutation<UpdateArticleResponse, UpdateArticleRequest>(UpdateArticle)

  // Small tips to prevent update warnings on unmounted components
  useEffect(
    () => () => {
      isMounted.current = false
    },
    []
  )

  const updateArticleStatus = useCallback(async () => {
    setLoading(true)
    try {
      await updateArticleMutation({
        variables: { id: article.id, status },
        update: updateCacheAfterUpdate,
      })
      if (onSuccess) onSuccess(article)
    } catch (err) {
      showErrorMessage(getGQLError(err))
    } finally {
      if (isMounted.current) {
        setLoading(false)
      }
    }
  }, [updateArticleMutation, article, status, onSuccess, showErrorMessage])

  const { title, icon, kbs } = variants[status]

  // Keyboard shortcut is only active for Floating Action Button
  useKeyboard(variants[status].kbs, updateArticleStatus, keyboard)
  const kbsLabel = keyboard ? ` [${kbs}]` : ''

  return (
    <ButtonIcon
      title={title + kbsLabel}
      onClick={updateArticleStatus}
      loading={loading}
      floating={floating}
      icon={icon}
      variant={status === 'to_read' ? 'default' : 'primary'}
    />
  )
}
