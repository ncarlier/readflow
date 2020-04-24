import { DataProxy } from 'apollo-cache'

import { GetCategoriesResponse } from '../categories/models'
import { GetCategories } from '../categories/queries'
import { AddNewArticleResponse, UpdateArticleResponse } from './models'
import { GetArticle } from './queries'

export const updateCacheAfterCreate = (proxy: DataProxy, mutationResult: { data?: AddNewArticleResponse | null }) => {
  if (!mutationResult || !mutationResult.data) {
    return
  }
  const article = mutationResult.data.addArticle
  // Update categories `_all` value
  try {
    const previousData = proxy.readQuery<GetCategoriesResponse>({
      query: GetCategories,
    })
    if (previousData && previousData.categories) {
      const { categories } = previousData
      categories._all++
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
  // Update categories `_all` and `_starred` values
  try {
    const previousData = proxy.readQuery<GetCategoriesResponse>({
      query: GetCategories,
    })
    if (previousData && previousData.categories) {
      const { categories } = previousData
      categories._all = updated._all
      categories._starred = updated._starred
      proxy.writeQuery({ data: { categories }, query: GetCategories })
    }
  } catch (err) {
    console.warn('unable to update categories cache when updating article')
  }
}
