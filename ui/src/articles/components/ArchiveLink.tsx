import React, { useCallback, useContext, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'

import { getGQLError } from '../../helpers'
import Kbd from '../../components/Kbd'
import LinkIcon from '../../components/LinkIcon'
import Loader from '../../components/Loader'
import { MessageContext } from '../../context/MessageContext'
import { ArchiveService } from '../../settings/archive-services/models'
import { Article } from '../models'
import { ArchiveArticle } from '../queries'

interface ArchiveArticleFields {
  id: number
  archiver: string
  noShortcuts?: boolean
}

interface Props {
  service: ArchiveService
  article: Article
  keyboard?: boolean
}

export default (props: Props) => {
  const [loading, setLoading] = useState(false)
  const { article, service, keyboard = false } = props
  const { showMessage, showErrorMessage } = useContext(MessageContext)

  const archiveArticleMutation = useMutation<ArchiveArticleFields>(ArchiveArticle)

  const archiveArticle = useCallback(async () => {
    setLoading(true)
    try {
      await archiveArticleMutation({
        variables: { id: article.id, archiver: service.alias }
      })
      showMessage(`Article sent to ${service.alias}: ${article.title}`)
    } catch (err) {
      showErrorMessage(getGQLError(err))
    } finally {
      setLoading(false)
    }
  }, [article])

  if (loading) {
    return <Loader />
  }

  return (
    <LinkIcon title={`Save to ${service.alias}`} icon="backup" onClick={archiveArticle}>
      <span>Save to {service.alias}</span>
      {keyboard && service.is_default && <Kbd keys="s" onKeypress={archiveArticle} />}
    </LinkIcon>
  )
}
