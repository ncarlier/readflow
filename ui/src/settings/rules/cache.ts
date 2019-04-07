import { DataProxy } from "apollo-cache"
import { CreateOrUpdateRuleResponse, GetRuleResponse, GetRulesResponse } from "./models";
import { GetRule, GetRules } from "./queries";

export const updateCacheAfterCreate = (proxy: DataProxy, mutationResult: {data: CreateOrUpdateRuleResponse}) => {
  const created = mutationResult.data.createOrUpdateRule
  const previousData = proxy.readQuery<GetRulesResponse>({
    query: GetRules,
  })
  previousData!.rules.unshift(created)
  proxy.writeQuery({ data: previousData, query: GetRules })
}

export const updateCacheAfterUpdate = (proxy: DataProxy, mutationResult: {data: CreateOrUpdateRuleResponse}) => {
  const updated = mutationResult!.data.createOrUpdateRule
  const previousData = proxy.readQuery<GetRulesResponse>({
    query: GetRules,
  })
  const rules = previousData!.rules.map(service => {
    return service.id === updated.id ? updated : service
  })
  proxy.writeQuery({ data: {rules}, query: GetRules })
  proxy.writeQuery({
    data: {
      rule: updated
    }, 
    query: GetRule,
    variables: {id: updated.id}
  })
}

export const updateCacheAfterDelete = (ids: number[]) => (proxy: DataProxy) => {
  const previousData = proxy.readQuery<GetRulesResponse>({
    query: GetRules,
  })
  const rules = previousData!.rules.filter(rule => !ids.includes(rule.id!))
  proxy.writeQuery({ data: {rules}, query: GetRules })
}
