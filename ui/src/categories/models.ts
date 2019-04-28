export interface Category {
  id?: number
  title: string
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

export interface DeleteCategoriesResponse {
  deleteCategories: number
}
