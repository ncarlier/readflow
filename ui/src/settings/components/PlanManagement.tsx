import React from 'react'

import { Button } from '../../components'
import { PORTAL_URL } from '../../config'
import { User } from '../../contexts'

interface Props {
  user: User
}

const UpgradePlanButton = () => (
  <Button as={'a'} href={`${PORTAL_URL}/pricing`} target="_blank" variant="primary" title="Upgrade your plan">
    Upgrade your plan
  </Button>
)

const ManagePlanButton = () => (
  <Button as={'a'} href={`${PORTAL_URL}/account`} target="_blank" title="Manage your plan">
    Manage your plan
  </Button>
)

export const PlanManagement = ({ user }: Props) => {
  if (PORTAL_URL === 'https://about.readflow.app') {
    return user.customer_id || user.plan === 'premium' ? <ManagePlanButton /> : <UpgradePlanButton />
  }
  return <p>Ask your administrator to update your plan if needed.</p>
}
