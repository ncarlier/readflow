import { DataProxy } from '@apollo/client'
import { GetCategoriesResponse } from '../categories/models'
import { GetCategories } from '../categories/queries'
import { AddNewArticleResponse, UpdateArticleResponse, MarkAllArticlesAsReadResponse } from './models'
import { GetArticle } from './queries'

export const updateCacheAfterCreate = (proxy: DataProxy, mutationResult: { data?: AddNewArticleResponse | null }) => {
  if (!mutationResult || !mutationResult.data) {
    return
  }
  const article = mutationResult.data.addArticle
  // Update categories `_inbox` value
  try {
    const previousData = proxy.readQuery<GetCategoriesResponse>({
      query: GetCategories,
    })
    if (previousData && previousData.categories) {
      const categories = { ...previousData.categories }
      categories._inbox++
      proxy.writeQuery({ data: { categories }, query: GetCategories })
    }
  } catch (err) {
    console.warn('unable to update categories cache when creating article')
  }
  // Update article
  proxy.writeQuery({
    data: { article },
    query: GetArticle,
    variables: { id: article.id },
  })
}

export const updateCacheAfterUpdate = (proxy: DataProxy, mutationResult: { data?: UpdateArticleResponse | null }) => {
  if (!mutationResult || !mutationResult.data) {
    return
  }
  const updated = mutationResult.data.updateArticle
  // Update categories `_inbox` and `_starred` values
  try {
    const previousData = proxy.readQuery<GetCategoriesResponse>({
      query: GetCategories,
    })
    if (previousData && previousData.categories) {
      const categories = { ...previousData.categories }
      categories._inbox = updated._inbox
      categories._to_read = updated._to_read
      categories._starred = updated._starred
      proxy.writeQuery({ data: { categories }, query: GetCategories })
    }
  } catch (err) {
    console.warn('unable to update categories cache when updating article')
  }
}

export const updateCacheAfterMarkAllAsRead = (
  proxy: DataProxy,
  mutationResult: { data?: MarkAllArticlesAsReadResponse | null }
) => {
  if (!mutationResult || !mutationResult.data) {
    return
  }
  const updated = mutationResult.data.markAllArticlesAsRead
  // Update categories and `_inbox` value
  try {
    const previousData = proxy.readQuery<GetCategoriesResponse>({
      query: GetCategories,
    })
    if (previousData && previousData.categories) {
      const categories = { ...previousData.categories }
      categories._inbox = updated._inbox
      const { entries } = updated
      // Merge categories inbox values
      categories.entries.forEach((cat) => {
        const found = entries.find((c) => cat.id === c.id)
        if (found) {
          Object.assign(cat, found)
        }
      })
      proxy.writeQuery({ data: { categories }, query: GetCategories })
    }
  } catch (err) {
    console.warn('unable to update categories cache when mark all as read')
  }
}
