import React from 'react'

import { Article } from '../../models'
import OfflineLink from './OfflineLink'
import ShareLink from './ShareLink'
import DownloadAsLink from './DownloadAsLink'
import OutgoingWebhooksMenuItems from './OutgoingWebhooksMenuItems'
import { DrawerMenu } from '../../../components'
import { EditArticleLink } from './EditArticleLink'

interface Props {
  article: Article
  keyboard?: boolean
  showEditModal?: () => void
}

export const ArticleContextMenu = (props: Props) => {
  const nvg: any = window.navigator
  const title = 'More actions ...'
  const isOnline = !props.article.isOffline
  
  return (
    <DrawerMenu title={title} kbs={props.keyboard ? 'm' : ''}>
      <ul>
        {props.showEditModal && (
          <li>
            <EditArticleLink {...props} showEditModal={props.showEditModal} />
          </li>
        )}
        {nvg.share && (
          <li>
            <ShareLink {...props} />
          </li>
        )}
        <li>
          <OfflineLink {...props} />
        </li>
        <li>
          <DownloadAsLink {...props} />
        </li>
        {isOnline && <OutgoingWebhooksMenuItems {...props} />}
      </ul>
    </DrawerMenu>
  )
}
