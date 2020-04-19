import { Category } from '../categories/models'

export interface Article {
  id: number
  title: string
  html: string
  text: string
  image: string
  url: string
  status: ArticleStatus
  starred: boolean
  category?: Category
  isOffline?: boolean
  published_at: string
  created_at: string
  updated_at: string
}

export type ArticleStatus = 'read' | 'unread'
export type SortOrder = 'asc' | 'desc'

export interface GetArticlesRequest {
  status: ArticleStatus | null
  starred: boolean | null
  category: number | null
  limit: number | null
  afterCursor: number | null
  sortOrder: SortOrder | null
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
  status?: ArticleStatus
  starred?: boolean
}

export interface UpdateArticleResponse {
  updateArticle: {
    article: Article
    _all: number
    _starred: number
  }
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
