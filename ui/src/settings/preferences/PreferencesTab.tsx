import React from 'react'

import Panel from '../../components/Panel'
import { usePageTitle } from '../../hooks'
import CleanHistoryBox from './CleanHistoryBox'
import DeleteAccountBox from './DeleteAccountBox'
import InstallationBox from './InstallationBox'
import NotificationBox from './NotificationBox'
import classes from './PreferencesTab.module.css'
import UserPlanSection from './UserPlanSection'

export default () => {
  usePageTitle('Settings - Preferences')

  return (
    <Panel className={classes.preferences}>
      <section>
        <h2>Device</h2>
        <hr />
        <p>Manage preferences on this device.</p>
        <InstallationBox />
        <NotificationBox />
        <UserPlanSection />
        <h2>Personal data</h2>
        <hr />
        <p>Data are yours and you have full control over it.</p>
        <CleanHistoryBox />
        <DeleteAccountBox />
      </section>
    </Panel>
  )
}
