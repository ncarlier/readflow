

export type Category = {
  id?: number
  title: string
  created_at?: string
  updated_at?: string
}

export type GetCategoriesResponse = {
  categories: Category[]
}

export interface GetCategoryResponse {
  category: Category
}

export type CreateOrUpdateCategoryResponse = {
  createOrUpdateCategory: Category
}

export type DeleteCategoriesResponse = {
  deleteCategories: number
}


