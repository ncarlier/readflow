import React from 'react'
import { Link, RouteComponentProps, withRouter } from 'react-router-dom'

import Icon from '../../components/Icon'
import TimeAgo from '../../components/TimeAgo'
import { classNames, getHostname } from '../../helpers'
import useKeyboard from '../../hooks/useKeyboard'
import { DropDownOrigin } from '../../components/DropdownMenu'
import { Article } from '../models'
import styles from './ArticleCard.module.css'
import ArticleImage from './ArticleImage'
import ArticleContextMenu from './context-menu/ArticleContextMenu'
import MarkAsButton from './MarkAsButton'
import StarsButton from './StarsButton'

interface Props {
  article: Article
  isActive: boolean
  onRemove?: () => void
}

const dropDownMenuOrigin: DropDownOrigin = { horizontal: 'left', vertical: 'top' }

type AllProps = Props & RouteComponentProps

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
        {article.url !== '' && (
          <a href={article.url} target="_blank" rel="noopener noreferrer" title="Open original article">
            <Icon name="open_in_new" />
            {getHostname(article.url)}
          </a>
        )}
        <TimeAgo dateTime={article.created_at} />
        <ArticleContextMenu article={article} keyboard={isActive} origin={dropDownMenuOrigin} />
        {!article.isOffline && <StarsButton article={article} keyboard={isActive} origin={dropDownMenuOrigin} />}
        {!article.isOffline && <MarkAsButton article={article} onSuccess={onRemove} keyboard={isActive} />}
      </footer>
    </article>
  )
})
