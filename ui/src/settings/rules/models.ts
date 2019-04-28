export interface Rule {
  id?: number
  alias: string
  category_id: number
  rule: string
  priority: number
  created_at?: string
  updated_at?: string
}

export interface GetRulesResponse {
  rules: Rule[]
}

export interface GetRuleResponse {
  rule: Rule
}

export interface CreateOrUpdateRuleResponse {
  createOrUpdateRule: Rule
}
