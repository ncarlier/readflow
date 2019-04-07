

export type Rule = {
  id?: number
  alias: string
  category_id: number
  rule: string
  priority: number
  created_at?: string
  updated_at?: string
}

export type GetRulesResponse = {
  rules: Rule[]
}

export interface GetRuleResponse {
  rule: Rule
}

export type CreateOrUpdateRuleResponse = {
  createOrUpdateRule: Rule
}
