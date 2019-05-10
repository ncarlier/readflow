import React from 'react'
import { Link, RouteComponentProps, withRouter } from 'react-router-dom'

import Panel from '../../common/Panel'
import useKeyboard from '../../hooks/useKeyboard'
import ArticleHeader from '../components/ArticleHeader'
import { Article } from '../models'
import styles from './ArticleCard.module.css'
import ArticleFooter from './ArticleFooter'
import ArticleMenu from './ArticleMenu'
import MarkAsButton from './MarkAsButton'

interface Props {
  article: Article
  isActive: boolean
  onRemove?: () => void
  readMoreBasePath: string
}

type AllProps = Props & RouteComponentProps

export default withRouter((props: AllProps) => {
  const { article, readMoreBasePath, isActive, onRemove, history } = props

  const readMorePath = readMoreBasePath + props.article.id

  useKeyboard('enter', () => history.push(readMorePath), isActive)
  const kbs = isActive ? ' [enter]' : ''

  return (
    <Panel className={isActive ? styles.active : ''}>
      <ArticleHeader article={article} to={readMorePath}>
        <ArticleMenu article={article} keyboard={isActive} />
      </ArticleHeader>
      <article className={styles.summary}>
        {article.image && (
          <Link to={readMorePath} className={styles.illustration} title={'Open article details' + kbs}>
            <img src={article.image} alt="Illustration" />
          </Link>
        )}
        {article.text && <p>{article.text}</p>}
      </article>
      {!article.isOffline && (
        <ArticleFooter>
          <MarkAsButton article={article} onSuccess={onRemove} keyboard={isActive} />
        </ArticleFooter>
      )}
    </Panel>
  )
})
