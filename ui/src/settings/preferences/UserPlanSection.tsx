import gql from 'graphql-tag'
import React from 'react'
import { useQuery } from '@apollo/client'

import Box from '../../components/Box'
import Loader from '../../components/Loader'
import { GetCurrentUser, GetCurrentUserResponse } from '../../components/UserInfos'
import ErrorPanel from '../../error/ErrorPanel'
import { matchResponse } from '../../helpers'
import PlanManagement from '../components/PlanManagement'

export const GetPlans = gql`
  query {
    plans {
      name
      total_articles
      total_categories
    }
  }
`
interface Plan {
  name: string
  total_articles: number
  total_categories: number
}

export interface GetPlansResponse {
  plans: Plan[]
}

interface UserPlanBoxProps {
  plans: Plan[]
}

const UserPlanBox = ({ plans }: UserPlanBoxProps) => {
  const { data, error, loading } = useQuery<GetCurrentUserResponse>(GetCurrentUser)

  const render = matchResponse<GetCurrentUserResponse>({
    Loading: () => <Loader center />,
    Error: (err) => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: ({ me }) => {
      let plan = plans.find((p) => p.name === me.plan)
      if (!plan) {
        plan = plans[0]
      }
      return (
        <Box title={plan.name}>
          <ul>
            <li>
              Max number of articles: <b>{plan.total_articles}</b>
            </li>
            <li>
              Max number of categories: <b>{plan.total_categories}</b>
            </li>
          </ul>
          <PlanManagement user={me} />
        </Box>
      )
    },
  })
  return <>{render(loading, data, error)}</>
}

export default () => {
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
