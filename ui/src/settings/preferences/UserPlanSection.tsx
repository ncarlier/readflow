import gql from 'graphql-tag'
import React from 'react'
import { useQuery } from 'react-apollo-hooks'

import Box from '../../components/Box'
import Loader from '../../components/Loader'
import { GetCurrentUser, GetCurrentUserResponse } from '../../components/UserInfos'
import ErrorPanel from '../../error/ErrorPanel'
import { matchResponse } from '../../helpers'

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
    Loading: () => <Loader />,
    Error: (err) => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: (data) => {
      let plan = plans.find((p) => p.name === data.me.plan)
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
          <p>You can ask administrator to update your plan if needed.</p>
        </Box>
      )
    },
    Other: () => <ErrorPanel>Unable to fetch current user plan!</ErrorPanel>,
  })
  return <>{render(data, error, loading)}</>
}

export default () => {
  const { data, error, loading } = useQuery<GetPlansResponse>(GetPlans)

  const render = matchResponse<GetPlansResponse>({
    Loading: () => <Loader />,
    Error: (err) => <ErrorPanel>{err.message}</ErrorPanel>,
    Data: (data) => {
      if (data.plans.length === 0) {
        return
      }
      return (
        <>
          <h2>User plan</h2>
          <hr />
          <p>Your user plan defines quotas and usage limits.</p>
          <UserPlanBox plans={data.plans} />
        </>
      )
    },
    Other: () => <ErrorPanel>Unable to fetch user plans!</ErrorPanel>,
  })

  return <>{render(data, error, loading)}</>
}
