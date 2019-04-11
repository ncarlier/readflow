import React from 'react'

import { Article } from '../models'

import styles from './ArticleFooter.module.css'
import MarkAsButton from './MarkAsButton'

type Props = {
  article: Article
}

type AllProps = Props

export default ({ article }: AllProps) => (
  <footer>
    <div className={styles.actions}>
      <MarkAsButton article={article} />
    </div>
  </footer>
)
