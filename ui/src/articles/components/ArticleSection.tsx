import React from 'react'

import { Article } from '../models'
import { MarkAsButton, StarsButton, ArticleContextMenu, ArticleHeader, ArticleContent } from '.'
import { useHistory } from 'react-router-dom'
import { useArticleEditModal } from '../hooks'

interface Props {
  article: Article
}

export const ArticleSection = ({article}: Props) => {
  const { goBack } = useHistory()
  const [showEditModal] = useArticleEditModal(article)
  return (
    <>
      <ArticleHeader article={article}>
        <StarsButton article={article} keyboard />
        {article.status === 'inbox' && (
          <MarkAsButton article={article} onSuccess={goBack} status="to_read" keyboard />
        )}
        <ArticleContextMenu article={article} keyboard showEditModal={showEditModal}/>
      </ArticleHeader>
      <ArticleContent article={article} />
      <MarkAsButton
        article={article}
        status={article.status === 'read' ? 'inbox' : 'read'}
        floating
        onSuccess={goBack}
        keyboard
      />
    </>
  )
}
