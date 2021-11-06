export interface Category {
  id?: number
  title: string
  rule: string | null
  inbox?: number
  created_at?: string
  updated_at?: string
}

export interface GetCategoriesResponse {
  categories: {
    _inbox: number
    _to_read: number
    _starred: number
    entries: Category[]
  }
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
