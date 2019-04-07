import React from 'react'

import { Article } from '../models'

import styles from './ArticleHeader.module.css'
import TimeAgo from '../../common/TimeAgo'
import Icon from '../../common/Icon'
import OfflineButton from './OfflineButton'
import MarkAsButton from './MarkAsButton'
import ArchiveButton from './ArchiveButton'

type Props = {
  article: Article
  showAllActions?: boolean
}

type AllProps = Props

export default ({article, showAllActions = false}: AllProps) => {
  const showPutOffline = showAllActions && !article.isOffline
  const showRemoveOffline = article.isOffline
  const showArchive = showAllActions
  const showMarkAs = !article.isOffline
  return (
    <header className={styles.header}>
      <h1>
        <small>
          <TimeAgo dateTime={article.created_at} />
          { article.url != "" &&
            <a href={article.url}
               target="_blank"
               rel="noopener noreferrer"
               title="Open original article">
              {(new URL(article.url)).hostname}
              <Icon name="open_in_new" />
            </a>
          }
        </small>
        <span>{article.title}</span>
      </h1>
      <div className={styles.actions}>
        { showPutOffline && <OfflineButton article={article} /> }
        { showArchive && <ArchiveButton article={article} /> }
        { showRemoveOffline && <OfflineButton article={article} remove /> }
        { showMarkAs && <MarkAsButton article={article} /> }
      </div>
    </header>
  )
}
