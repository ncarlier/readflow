import React from 'react'
import { Link, RouteComponentProps, withRouter } from 'react-router-dom'

import useKeyboard from '../../hooks/useKeyboard'
import { Article } from '../models'
import styles from './ArticleCard.module.css'
import ArticleMenu from './ArticleMenu'
import Icon from '../../common/Icon'
import TimeAgo from '../../common/TimeAgo'
import { classNames } from '../../common/helpers';
import MarkAsButton from './MarkAsButton';

interface Props {
  article: Article
  isActive: boolean
  onRemove?: () => void
}

type AllProps = Props & RouteComponentProps

const menuStyle = {
  top: 'initial',
  right: 0,
  bottom: '100%'
}

export default withRouter((props: AllProps) => {
  const { article, isActive, onRemove, history, match } = props

  const readMorePath = match.url + '/' + props.article.id

  useKeyboard('enter', () => history.push(readMorePath), isActive)
  const kbs = isActive ? ' [enter]' : ''
  const activeClass = isActive ? styles.active : ''

  return (
    <article className={classNames(styles.card, activeClass)}>
      { article.image && (
        <Link to={readMorePath} title="View details" className={styles.illustration}>
          <img src={article.image} alt="Illustration" onError={(e) => e.currentTarget.classList.add(styles.broken)} />
        </Link>
      )}
      <Link to={readMorePath} title="View details" className={styles.content}>
        <header>
          {article.title}
        </header>
        {article.text && <p>{article.text}</p>}
      </Link>
      <footer>
        {article.url != '' && (
          <a href={article.url} target="_blank" rel="noopener noreferrer" title="Open original article">
            {new URL(article.url).hostname}
            <Icon name="open_in_new" />
          </a>
        )}
        <TimeAgo dateTime={article.created_at} />
        <ArticleMenu article={article} keyboard={isActive} style={menuStyle}/>
        {!article.isOffline && <MarkAsButton article={article} onSuccess={onRemove} keyboard={isActive} />}
      </footer>
    </article>
  )
})
