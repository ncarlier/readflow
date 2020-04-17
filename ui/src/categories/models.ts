export interface Category {
  id?: number
  title: string
  rule: string | null
  unread?: number
  created_at?: string
  updated_at?: string
}

export interface GetCategoriesResponse {
  categories: Category[]
}

export interface GetCategoryResponse {
  category: Category
}

export interface CreateOrUpdateCategoryResponse {
  createOrUpdateCategory: Category
}

export interface DeleteCategoriesRequest {
  ids: number[]
}

export interface DeleteCategoriesResponse {
  deleteCategories: number
}
