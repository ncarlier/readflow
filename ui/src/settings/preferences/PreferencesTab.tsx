import React from 'react'

import Panel from '../../components/Panel'
import { usePageTitle } from '../../hooks'
import InstallationBox from './InstallationBox'
import NotificationBox from './NotificationBox'
import classes from './PreferencesTab.module.css'

export default () => {
  usePageTitle('Settings - Preferences')

  return (
    <Panel className={classes.preferences}>
      <section>
        <h2>Device</h2>
        <hr />
        <p>Preferences on this device.</p>
        <InstallationBox />
        <NotificationBox />
      </section>
    </Panel>
  )
}
