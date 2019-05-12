import React, { useCallback, useState } from 'react'
import { useMutation } from 'react-apollo-hooks'

import { getGQLError } from '../../common/helpers'
import Kbd from '../../common/Kbd'
import LinkIcon from '../../common/LinkIcon'
import Loader from '../../common/Loader'
import { connectMessageDispatch, IMessageDispatchProps } from '../../containers/MessageContainer'
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

type AllProps = Props & IMessageDispatchProps

export const ArchiveLink = (props: AllProps) => {
  const [loading, setLoading] = useState(false)
  const { article, service, showMessage, keyboard = false } = props

  const archiveArticleMutation = useMutation<ArchiveArticleFields>(ArchiveArticle)

  const archiveArticle = useCallback(async () => {
    setLoading(true)
    try {
      await archiveArticleMutation({
        variables: { id: article.id, archiver: service.alias }
      })
      showMessage(`Article sent to ${service.alias}: ${article.title}`)
    } catch (err) {
      showMessage(getGQLError(err), true)
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

export default connectMessageDispatch(ArchiveLink)
