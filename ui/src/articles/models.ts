import { Category } from '../categories/models'

export interface Article {
  id: number
  title: string
  html: string
  text: string
  image: string
  url: string
  status: string
  category?: Category
  isOffline?: boolean
  published_at: string
  created_at: string
  updated_at: string
}

export interface GetArticlesRequest {
  status: string
  category?: number
  limit: number
  afterCursor?: number
  sortOrder: string
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

export interface UpdateArticleStatusRequest {
  id: number
  status: string
}

export interface UpdateArticleStatusResponse {
  updateArticleStatus: Article
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
