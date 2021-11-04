import React from 'react'
import { Link, RouteComponentProps, withRouter } from 'react-router-dom'

import { Icon, TimeAgo } from '../../components'
import { classNames, getHostname } from '../../helpers'
import { useKeyboard } from '../../hooks'
import { Article } from '../models'
import styles from './ArticleCard.module.css'
import { ArticleImage, MarkAsButton, StarsButton, ArticleContextMenu } from '.'

interface Props {
  article: Article
  isActive: boolean
  onRemove?: () => void
}

type AllProps = Props & RouteComponentProps

export const ArticleCard = withRouter((props: AllProps) => {
  const { article, isActive, onRemove, history, match } = props

  const readMorePath = match.url + '/' + props.article.id

  useKeyboard('enter', () => history.push(readMorePath), isActive)
  const kbs = isActive ? ' [enter]' : ''
  const activeClass = isActive ? styles.active : ''

  return (
    <article className={classNames(styles.card, activeClass)}>
      {article.image && (
        <Link to={readMorePath} title={'View details' + kbs} className={styles.illustration}>
          <ArticleImage src={article.image} alt={article.title} />
        </Link>
      )}
      <Link to={readMorePath} title={'View details' + kbs} className={styles.content}>
        {article.category && <h3>{article.category.title}</h3>}
        <header>{article.title}</header>
        {article.text && <p>{article.text}</p>}
      </Link>
      <footer>
        {article.url !== '' && (
          <a href={article.url} target="_blank" rel="noopener noreferrer" title="Open original article">
            <Icon name="open_in_new" />
            {getHostname(article.url)}
          </a>
        )}
        <TimeAgo dateTime={article.created_at} />
        <ArticleContextMenu article={article} keyboard={isActive} />
        {!article.isOffline && (
          <>
            {article.status != 'inbox' && <StarsButton article={article} keyboard={isActive} />}
            {article.status === 'inbox' && (
              <MarkAsButton article={article} status="to_read" onSuccess={onRemove} keyboard={isActive} />
            )}
            {article.status != 'read' && (
              <MarkAsButton article={article} status="read" onSuccess={onRemove} keyboard={isActive} />
            )}
            {article.status === 'read' && (
              <MarkAsButton article={article} status="inbox" onSuccess={onRemove} keyboard={isActive} />
            )}
          </>
        )}
      </footer>
    </article>
  )
})
