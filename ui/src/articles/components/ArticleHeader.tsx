import { History } from 'history'
import React, { ReactNode } from 'react'
import { Link } from 'react-router-dom'

import Icon from '../../common/Icon'
import TimeAgo from '../../common/TimeAgo'
import { Article } from '../models'
import styles from './ArticleHeader.module.css'

interface Props {
  article: Article
  to?: History.LocationDescriptor
  children?: ReactNode
}

type AllProps = Props

export default ({ article, to, children }: AllProps) => (
  <header className={styles.header}>
    <h1>
      <small>
        <TimeAgo dateTime={article.created_at} />
        {article.url != '' && (
          <a href={article.url} target="_blank" rel="noopener noreferrer" title="Open original article">
            {new URL(article.url).hostname}
            <Icon name="open_in_new" />
          </a>
        )}
      </small>
      <span>
        {to ? (
          <Link to={to} title="View details">
            {article.title}
          </Link>
        ) : (
          article.title
        )}
      </span>
    </h1>
    <div className={styles.actions}>{children}</div>
  </header>
)
