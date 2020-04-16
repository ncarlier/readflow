import React from 'react'
import { Link, RouteComponentProps, withRouter } from 'react-router-dom'

import Icon from '../../components/Icon'
import TimeAgo from '../../components/TimeAgo'
import { classNames, getHostname } from '../../helpers'
import useKeyboard from '../../hooks/useKeyboard'
import { Article } from '../models'
import styles from './ArticleCard.module.css'
import ArticleImage from './ArticleImage'
import ArticleMenu from './ArticleMenu'
import MarkAsButton from './MarkAsButton'
import StarsButton from './StarsButton'

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
        {article.url != '' && (
          <a href={article.url} target="_blank" rel="noopener noreferrer" title="Open original article">
            {getHostname(article.url)}
            <Icon name="open_in_new" />
          </a>
        )}
        <TimeAgo dateTime={article.created_at} />
        <ArticleMenu article={article} keyboard={isActive} style={menuStyle} />
        {!article.isOffline && <StarsButton article={article} keyboard={isActive} />}
        {!article.isOffline && <MarkAsButton article={article} onSuccess={onRemove} keyboard={isActive} />}
      </footer>
    </article>
  )
})
