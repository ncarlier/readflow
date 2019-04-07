import React, { useState, useCallback } from 'react'
import { useMutation, useApolloClient } from 'react-apollo-hooks'

import { Article } from '../models'
import ButtonIcon from '../../common/ButtonIcon'

import { ArchiveArticle } from '../queries'
import { getGQLError } from '../../common/helpers'
import { GetArchiveServicesResponse } from '../../settings/archive-services/models'
import { GetArchiveServices } from '../../settings/archive-services/queries'
import { IMessageDispatchProps, connectMessageDispatch } from '../../containers/MessageContainer';
import useConfirmModal from '../../hooks/useConfirmModal'

type ArchiveArticleFields = {
  id: number
  archiver: string
}

type Props = {
  article: Article
}

type AllProps = Props & IMessageDispatchProps

export const ArchiveButton = (props: AllProps) => {
  const {
    article,
    showMessage
  } = props

  const client = useApolloClient()
  const [loading, setLoading] = useState(false)
  const archiveArticleMutation = useMutation<ArchiveArticleFields>(ArchiveArticle)
  const [showNoArchiveServiceModal] = useConfirmModal(
    'Archiving',
    <p>
      No archive service configured.<br/>
      Please configure one in the <a href="/settings/archive-services">setting page</a>.
    </p>
  )
  
  const archiveArticle = async (alias: string) => {
    await archiveArticleMutation({
      variables: {id: article.id, archiver: alias}
    })
  }

  const selectAndUseArchiveService = async () => {
    try {
      setLoading(true)
      const { errors, data } = await client.query<GetArchiveServicesResponse>({
        query: GetArchiveServices,
      })
      if (data) {
        if (data.archivers.length === 0) {
          showNoArchiveServiceModal()
        } else if (data.archivers.length > 1) {
          // TODO: Show choosing modal
          showMessage('You have to choose a default service archiver. Abort.')
        } else {
          await archiveArticle(data.archivers[0].alias)
          showMessage(`Article put offline: ${article.title}`)
        }
      }
      setLoading(false)
      if (errors) {
        throw new Error(errors[0])
      }
    } catch (err) {
      showMessage(getGQLError(err), true)
    }
  }    

  const handleOnClick = useCallback(() => {
    selectAndUseArchiveService()
  }, [article])

  return (
    <ButtonIcon
      title="Save to your cloud provider"
      icon="backup"
      onClick={handleOnClick}
      loading={loading}
    />
  )
}

export default connectMessageDispatch(ArchiveButton)
