import { Category } from '../categories/models'

export interface Article {
  id: number
  title: string
  html: string
  text: string
  image: string
  url: string
  status: string
  starred: boolean
  category?: Category
  isOffline?: boolean
  published_at: string
  created_at: string
  updated_at: string
}

export interface GetArticlesRequest {
  status: string | null
  starred: boolean | null
  category: number | null
  limit: number | null
  afterCursor: number | null
  sortOrder: string | null
}

export interface GetArticlesResponse {
  articles: {
    totalCount: number
    endCursor: number
    hasNext: boolean
    entries: Article[]
  }
}

export interface GetArticleResponse {
  article: Article
}

export interface UpdateArticleRequest {
  id: number
  status?: string
  starred?: string
}

export interface UpdateArticleResponse {
  updateArticle: Article
}

export interface MarkAllArticlesAsReadResponse {
  markAllArticlesAsRead: number
}

export interface AddNewArticleRequest {
  url: string
  category?: number
}

export interface AddNewArticleResponse {
  addArticle: Article
}
