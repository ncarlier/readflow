import React from 'react'

import { Article } from '../models'

import styles from './ArticleFooter.module.css'
import OfflineButton from './OfflineButton'
import ArchiveButton from './ArchiveButton'
import Button from '../../common/Button'

type Props = {
  article: Article
  readMorePath: string
}

type AllProps = Props

export default ({article, readMorePath}: AllProps) => (
  <footer>
    <div className={styles.content}>
      <Button
        title="Open article details"
        to={readMorePath}>
        Read more
      </Button>
    </div>
    <div className={styles.actions}>
      { !article.isOffline && <OfflineButton article={article} /> }
      <ArchiveButton article={article} />
    </div>
  </footer>
)
