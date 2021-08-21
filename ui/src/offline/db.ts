import Dexie from 'dexie'

import { Article } from '../articles/models'

class OfflineDatabase extends Dexie {
  articles: Dexie.Table<Article, number>

  constructor() {
    super('readflow-offline')
    this.version(1).stores({
      articles: '++id, title, created_at',
    })
    this.articles = this.table('articles')
  }
}

export const db = new OfflineDatabase()
