import React, { FC, PropsWithChildren } from 'react'
import { History } from 'history'
import { Link } from 'react-router-dom'

import { Icon, TimeAgo } from '../../components'
import { getHostname } from '../../helpers'
import { Article } from '../models'
import styles from './ArticleHeader.module.css'

interface Props extends PropsWithChildren {
  article: Article
  to?: History.LocationDescriptor
}

export const ArticleHeader: FC<Props> = ({ article, to, children }) => (
  <header className={styles.header}>
    <h1>
      {article.category && <small>{article.category.title}</small>}
      <span>
        {to ? (
          <Link to={to} title="View details">
            {article.title}
          </Link>
        ) : (
          article.title
        )}
      </span>
      <small>
        {article.url && (
          <a href={article.url} target="_blank" rel="noopener noreferrer" title="Open original article">
            {getHostname(article.url)}
            <Icon name="open_in_new" />
          </a>
        )}
        <TimeAgo dateTime={article.created_at} />
      </small>
    </h1>
    <div className={styles.actions}>{children}</div>
  </header>
)
