import React from 'react'

import Button from '../../components/Button'
import { User } from '../../components/UserInfos'
import { REDIRECT_URL } from '../../constants'

interface Props {
  user: User
}

const UpgradePlanButton = () => (
  <Button as={'a'} href={`${REDIRECT_URL}/pricing`} target="_blank" variant="primary" title="Upgrade your plan">
    Upgrade your plan
  </Button>
)

const ManagePlanButton = () => (
  <Button as={'a'} href={`${REDIRECT_URL}/account`} target="_blank" tile="Manage your plan">
    Manage your plan
  </Button>
)

export default ({ user }: Props) => {
  if (REDIRECT_URL === 'https://about.readflow.app') {
    return user.customer_id ? <ManagePlanButton /> : <UpgradePlanButton />
  }
  return <p>Ask your administrator to update your plan if needed.</p>
}
