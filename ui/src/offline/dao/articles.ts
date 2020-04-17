import { Article, GetArticlesRequest, GetArticlesResponse } from '../../articles/models'
import { db } from '../db'

export const saveArticle = (article: Article) => {
  return db.transaction('rw', db.articles, async () => {
    article.isOffline = true
    const id = await db.articles.put(article)
    console.log('Article put into offline storage:', id)
    return article
  })
}

export const removeArticle = (article: Article) => {
  return db.transaction('rw', db.articles, async () => {
    const id = await db.articles.delete(article.id)
    console.log('Article removed from offline storage:', id)
    return article
  })
}

export const getArticle = (id: number) => {
  return db.articles.get(id)
}

export const getTotalNbArticles = () => {
  db.articles.count()
}

export const getArticles = async (req: GetArticlesRequest) => {
  const { afterCursor, sortOrder = 'asc' } = req
  const limit = req.limit ? req.limit : 10
  const table = db.articles

  const result: GetArticlesResponse = {
    articles: {
      endCursor: -1,
      entries: [],
      hasNext: false,
      totalCount: 0
    }
  }
  result.articles.totalCount = await table.count()

  const asc = sortOrder === 'asc'
  if (afterCursor) {
    let collection = table.orderBy('id')
    if (!asc) {
      collection = collection.reverse()
    }
    const pageKeys: number[] = []
    await collection
      .until(() => pageKeys.length === limit + 1)
      .eachPrimaryKey(id => {
        if ((asc && id > afterCursor) || (!asc && id < afterCursor)) {
          pageKeys.push(id)
        }
      })
    result.articles.entries = await Promise.all<Article>(pageKeys.map(id => table.get(id)))
  } else {
    let collection = table.orderBy('id')
    if (!asc) {
      collection = collection.reverse()
    }
    result.articles.entries = await collection.limit(limit + 1).toArray()
  }

  if (result.articles.entries.length > limit) {
    result.articles.entries.pop()
    result.articles.hasNext = true
  }

  if (result.articles.entries.length) {
    result.articles.endCursor = (result.articles.entries[result.articles.entries.length - 1] as Article).id
  }

  return result
}
