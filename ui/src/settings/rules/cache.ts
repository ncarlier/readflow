import { DataProxy } from 'apollo-cache'

import { CreateOrUpdateRuleResponse, GetRulesResponse } from './models'
import { GetRule, GetRules } from './queries'

export const updateCacheAfterCreate = (
  proxy: DataProxy,
  mutationResult: { data?: CreateOrUpdateRuleResponse | null }
) => {
  if (!mutationResult || !mutationResult.data) {
    return
  }
  const created = mutationResult.data.createOrUpdateRule
  const previousData = proxy.readQuery<GetRulesResponse>({
    query: GetRules
  })
  if (previousData) {
    previousData.rules.unshift(created)
    proxy.writeQuery({ data: previousData, query: GetRules })
  }
}

export const updateCacheAfterUpdate = (
  proxy: DataProxy,
  mutationResult: { data?: CreateOrUpdateRuleResponse | null }
) => {
  if (!mutationResult || !mutationResult.data) {
    return
  }
  const updated = mutationResult.data.createOrUpdateRule
  const previousData = proxy.readQuery<GetRulesResponse>({
    query: GetRules
  })
  if (previousData) {
    const rules = previousData.rules.map(service => {
      return service.id === updated.id ? updated : service
    })
    proxy.writeQuery({ data: { rules }, query: GetRules })
  }
  proxy.writeQuery({
    data: {
      rule: updated
    },
    query: GetRule,
    variables: { id: updated.id }
  })
}

export const updateCacheAfterDelete = (ids: number[]) => (proxy: DataProxy) => {
  const previousData = proxy.readQuery<GetRulesResponse>({
    query: GetRules
  })
  if (previousData) {
    const rules = previousData.rules.filter(rule => rule.id && !ids.includes(rule.id))
    proxy.writeQuery({ data: { rules }, query: GetRules })
  }
}
