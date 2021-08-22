import React from 'react'

import { Article } from '../../models'
import OfflineLink from './OfflineLink'
import ShareLink from './ShareLink'
import DownloadAsLink from './DownloadAsLink'
import OutgoingWebhooksMenuItems from './OutgoingWebhooksMenuItems'
import { DrawerMenu } from '../../../components'

interface Props {
  article: Article
  keyboard?: boolean
}

export const ArticleContextMenu = (props: Props) => {
  const nvg: any = window.navigator
  const title = 'More actions...'
  const isOnline = !props.article.isOffline

  return (
    <DrawerMenu title={title}>
      <ul>
        {nvg.share && (
          <li>
            <ShareLink {...props} />
          </li>
        )}
        <li>
          <DownloadAsLink {...props} />
        </li>
        <li>
          <OfflineLink {...props} />
        </li>
        {isOnline && <OutgoingWebhooksMenuItems {...props} />}
      </ul>
    </DrawerMenu>
  )
}
