import { Category } from '../categories/models'

export interface Article {
  id: number
  title: string
  html: string
  text: string
  image: string
  thumbhash: string
  thumbnails?: ArticleThumbnail[]
  url: string
  status: ArticleStatus
  stars: number
  category?: Category
  isOffline?: boolean
  published_at: string
  created_at: string
  updated_at: string
}

export interface ArticleThumbnail {
  size: string
  hash: string
}

export type ArticleStatus = 'inbox' | 'read' | 'to_read'
export type SortOrder = 'asc' | 'desc'
export type SortBy = 'key' | 'stars'

export interface GetArticlesRequest {
  status: ArticleStatus | null
  starred: boolean | null
  category: number | null
  limit: number | null
  afterCursor: number | null
  sortOrder: SortOrder | null
  sortBy: SortBy | null
  query: string | null
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
  title?: string
  text?: string
  status?: ArticleStatus
  stars?: number
  category_id?: number
  refresh?: boolean
}

export interface UpdateArticleResponse {
  updateArticle: {
    article: Article
    _inbox: number
    _to_read: number
    _starred: number
  }
}

export interface SendArticleToOutgoingWebhookResponse {
  sendArticleToOutgoingWebhook: {
    url?: string
    text?: string
    json?: string
  }
}

export interface MarkAllArticlesAsReadRequest {
  category: number | null
  status: ArticleStatus
}

export interface MarkAllArticlesAsReadResponse {
  markAllArticlesAsRead: {
    _inbox: number
    entries: Category[]
  }
}

export interface AddNewArticleRequest {
  url: string
  category?: number
}

export interface AddNewArticleResponse {
  addArticle: Article
}
