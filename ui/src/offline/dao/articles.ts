import { db } from "../db"
import { Article } from "../../articles/models"

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

export interface GetArticlesQuery {
  limit: number
  afterCursor?: number
  sortOrder: string
}

export interface GetArticlesResult {
  totalCount: number
  endCursor: number
  hasNext: boolean
  entries: Article[]
}

export const getTotalNbArticles = () => {
  db.articles.count()
}

export const getArticles = async ({limit = 10, afterCursor, sortOrder = 'asc'}: GetArticlesQuery) => {
  const table = db.articles
  
  const result: GetArticlesResult = {
    endCursor: -1,
    entries: [],
    hasNext: false,
    totalCount: 0,
  }
  result.totalCount = await table.count() 

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
    result.entries = await Promise.all<Article>(pageKeys.map(id => table.get(id)))
  } else {
    let collection = table.orderBy('id')
    if (!asc) {
      collection = collection.reverse()
    }
    result.entries = await collection.limit(limit + 1).toArray()
  }

  if (result.entries.length > limit) {
    result.entries.pop()
    result.hasNext = true
  }

  if (result.entries.length) {
    result.endCursor = (result.entries[result.entries.length - 1] as Article).id
  }
  
 return result
}
