import React from 'react'

import Panel from '../../components/Panel'
import { usePageTitle } from '../../hooks'
import CleanHistoryBox from './CleanHistoryBox'
import DeleteAccountBox from './DeleteAccountBox'
import InstallationBox from './InstallationBox'
import NotificationBox from './NotificationBox'
import ThemeBox from './ThemeBox'
import UserPlanSection from './UserPlanSection'

export default () => {
  usePageTitle('Settings - Preferences')

  return (
    <Panel>
      <section>
        <header>
          <h2>Device</h2>
        </header>
        <p>Manage preferences on this device.</p>
        <InstallationBox />
        <NotificationBox />
        <ThemeBox />
      </section>
      <UserPlanSection />
      <section>
        <header>
          <h2>Personal data</h2>
        </header>
        <p>Data are yours and you have full control over it.</p>
        <CleanHistoryBox />
        <DeleteAccountBox />
      </section>
    </Panel>
  )
}
