import gql from 'graphql-tag'
import React from 'react'
import { useQuery } from '@apollo/client'

import { Box, ErrorPanel, Loader } from '../../components'
import { matchResponse } from '../../helpers'
import { PlanManagement } from '../components'
import { useCurrentUser } from '../../contexts'

export const GetPlans = gql`
  query {
    plans {
      name
      articles_limit
      categories_limit
      incoming_webhooks_limit
      outgoing_webhooks_limit
    }
  }
`
interface Plan {
  name: string
  articles_limit: number
  categories_limit: number
  incoming_webhooks_limit: number
  outgoing_webhooks_limit: number
}

export interface GetPlansResponse {
  plans: Plan[]
}

interface UserPlanBoxProps {
  plans: Plan[]
}

const UserPlanBox = ({ plans }: UserPlanBoxProps) => {
  const user = useCurrentUser()

  if (!user) {
    return null
  }

  let plan = plans.find((p) => p.name === user.plan)
  if (!plan) {
    plan = plans[0]
  }
  return (
    <Box title={plan.name}>
      <ul>
        <li>
          Up to <b>{plan.articles_limit}</b> articles.
        </li>
        <li>
          Up to <b>{plan.categories_limit}</b> categories.
        </li>
        <li>
          Up to <b>{plan.incoming_webhooks_limit}</b> incoming webhooks.
        </li>
        <li>
          Up to <b>{plan.outgoing_webhooks_limit}</b> outgoing webhooks.
        </li>
        {plan.name === 'premium' && (
          <li>
            RSS feeds with a dedicated&nbsp;
            <a href={`https://feedpushr.nunux.org/${user.hashid}`} rel="noreferrer noopener" target="_blank">
              Feedpushr instance
            </a>
            .
          </li>
        )}
      </ul>
      <PlanManagement user={user} />
    </Box>
  )
}

const UserPlanSection = () => {
  const { data, error, loading } = useQuery<GetPlansResponse>(GetPlans)

  const render = matchResponse<GetPlansResponse>({
    Loading: () => <Loader center />,
    Error: (err) => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: (data) => {
      if (data.plans.length === 0) {
        return
      }
      return (
        <section>
          <header>
            <h2 id="plan">User plan</h2>
          </header>
          <p>Your user plan defines quotas and usage limits.</p>
          <UserPlanBox plans={data.plans} />
        </section>
      )
    },
  })

  return <>{render(loading, data, error)}</>
}

export default UserPlanSection
