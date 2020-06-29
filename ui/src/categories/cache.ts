import { DataProxy } from 'apollo-cache'

import { CreateOrUpdateCategoryResponse, GetCategoriesResponse } from './models'
import { GetCategories } from './queries'

export const updateCacheAfterCreate = (
  proxy: DataProxy,
  mutationResult: { data?: CreateOrUpdateCategoryResponse | null }
) => {
  if (!mutationResult.data) {
    return
  }
  const created = mutationResult.data.createOrUpdateCategory
  created.unread = 0
  const previousData = proxy.readQuery<GetCategoriesResponse>({
    query: GetCategories,
  })
  if (previousData) {
    const { categories } = previousData
    categories.entries = [created, ...categories.entries]
    proxy.writeQuery<GetCategoriesResponse>({ data: { categories }, query: GetCategories })
  }
}

export const updateCacheAfterDelete = (ids: number[]) => (proxy: DataProxy) => {
  const previousData = proxy.readQuery<GetCategoriesResponse>({
    query: GetCategories,
  })
  if (previousData) {
    const { categories } = previousData
    categories.entries = categories.entries.filter((category) => category.id && !ids.includes(category.id))
    proxy.writeQuery<GetCategoriesResponse>({ data: { categories }, query: GetCategories })
  }
}
