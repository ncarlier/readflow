import React, { useCallback } from 'react'
import { useMutation } from 'react-apollo-hooks'

import { getGQLError } from '../../common/helpers'
import LinkIcon from '../../common/LinkIcon'
import { connectMessageDispatch, IMessageDispatchProps } from '../../containers/MessageContainer'
import useKeyboard from '../../hooks/useKeyboard'
import { ArchiveService } from '../../settings/archive-services/models'
import { Article } from '../models'
import { ArchiveArticle } from '../queries'

interface ArchiveArticleFields {
  id: number
  archiver: string
  noShortcuts?: boolean
}

interface Props {
  article: Article
  service: ArchiveService
  noShortcuts?: boolean
}

type AllProps = Props & IMessageDispatchProps

export const ArchiveLink = (props: AllProps) => {
  const { article, service, showMessage, noShortcuts } = props

  const archiveArticleMutation = useMutation<ArchiveArticleFields>(ArchiveArticle)

  const archiveArticle = useCallback(async () => {
    try {
      await archiveArticleMutation({
        variables: { id: article.id, archiver: service.alias }
      })
      showMessage(`Article sent to ${service.alias}: ${article.title}`)
    } catch (err) {
      showMessage(getGQLError(err), true)
    }
  }, [article])

  useKeyboard('s', archiveArticle, service.is_default && !noShortcuts)

  return (
    <LinkIcon title={`Save to ${service.alias}`} icon="backup" onClick={archiveArticle}>
      <span>Save to {service.alias}</span>
      {service.is_default && !noShortcuts && <kbd>s</kbd>}
    </LinkIcon>
  )
}

export default connectMessageDispatch(ArchiveLink)
