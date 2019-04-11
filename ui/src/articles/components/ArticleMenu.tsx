import React from 'react'

import { Article } from '../models'

import DropdownMenu from '../../common/DropdownMenu'
import OfflineLink from './OfflineLink'
import ArchiveLink from './ArchiveLink'

type Props = {
  article: Article
}

type AllProps = Props

export default ({ article }: AllProps) => {
  return (
    <DropdownMenu>
      <ul>
        { article.isOffline ?
          <li>
            <OfflineLink article={article} remove />
          </li>
          :
          <li>
            <OfflineLink article={article} />
          </li>
        }
        <li>
          <ArchiveLink article={article} />
        </li>
      </ul>
    </DropdownMenu>
  )
}
