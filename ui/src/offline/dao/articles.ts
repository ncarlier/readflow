import { db } from "../db"
import { Article } from "../../articles/models"

export const saveArticle = (article: Article) => {
  return db.transaction('rw', db.articles, async () => {
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

export const getArticles = () => {
  return db.articles.toArray()
}
