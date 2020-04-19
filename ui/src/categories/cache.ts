import { DataProxy } from 'apollo-cache'

import { CreateOrUpdateCategoryResponse, GetCategoriesResponse } from './models'
import { GetCategories } from './queries'

export const updateCacheAfterCreate = (
  proxy: DataProxy,
  mutationResult: { data?: CreateOrUpdateCategoryResponse | null }
) => {
  const previousData = proxy.readQuery<GetCategoriesResponse>({
    query: GetCategories
  })
  if (previousData && mutationResult && mutationResult.data) {
    previousData.categories.entries.unshift(mutationResult.data.createOrUpdateCategory)
  }
  proxy.writeQuery({ data: previousData, query: GetCategories })
}
/*
export const updateCacheAfterUpdate = (
  proxy: DataProxy,
  mutationResult: { data?: CreateOrUpdateCategoryResponse | null }
) => {
  if (!mutationResult || !mutationResult.data) {
    return
  }
  const updated = mutationResult.data.createOrUpdateCategory
  const previousData = proxy.readQuery<GetCategoriesResponse>({
    query: GetCategories
  })
  if (previousData) {
    const categories = previousData.entries.map(cat => {
      return cat.id === updated.id ? updated : cat
    })
    proxy.writeQuery({ data: { categories }, query: GetCategories })
  }
  proxy.writeQuery({
    data: {
      category: updated
    },
    query: GetCategory,
    variables: { id: updated.id }
  })
}*/

export const updateCacheAfterDelete = (ids: number[]) => (proxy: DataProxy) => {
  const previousData = proxy.readQuery<GetCategoriesResponse>({
    query: GetCategories
  })
  if (previousData) {
    const categories = previousData.categories.entries.filter(category => category.id && !ids.includes(category.id))
    proxy.writeQuery({ data: { categories }, query: GetCategories })
  }
}
