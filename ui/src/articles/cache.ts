import { DataProxy } from "apollo-cache"
import { GetArticlesResponse, UpdateArticleStatusResponse, GetArticleResponse } from "./models"
import { GetArticle, GetArticles } from "./queries"


export const updateCacheAfterUpdateStatus = (proxy: DataProxy, mutationResult: {data: UpdateArticleStatusResponse}) => {
  const updated = mutationResult!.data.updateArticleStatus
  // Update GetArticle cache
  let previousArticleResponse: GetArticleResponse | null
  try {
    previousArticleResponse = proxy.readQuery<GetArticleResponse>({
      query: GetArticle,
      variables: {id: updated.id+""} // +"" is a hack because ID
    })
  } catch (e) {
    previousArticleResponse = null
  }
  if (previousArticleResponse) {
    const merged = Object.assign({}, previousArticleResponse.article, updated)
    proxy.writeQuery({
      data: {
        article: merged
      }, 
      query: GetArticle,
      variables: {id: updated.id}
    })
  }
  // Update GetArticles cache
  let previousArticlesResponse: GetArticlesResponse | null
  try {
    previousArticlesResponse = proxy.readQuery<GetArticlesResponse>({
      query: GetArticles,
    })
  } catch (e) {
    previousArticlesResponse = null
  }
  if (previousArticlesResponse) {
    previousArticlesResponse.articles.entries = previousArticlesResponse.articles.entries.map(article => {
      if (article.id === updated.id) {
        const merged = Object.assign({}, article, updated)
        return merged
      }
      return article
    })
    proxy.writeQuery({ data: previousArticlesResponse, query: GetArticles })
  }
}
