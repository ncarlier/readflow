import React from 'react'
import { Redirect, Route, RouteComponentProps, Switch } from 'react-router'

import { Appbar, Page } from '../layout'
import AboutButton from './about/AboutButton'
import CategoriesRoutes from './categories/routes'
import { Header, Tabs } from './components'
import IntegrationsRoutes from './intergrations/routes'
import PreferencesTab from './preferences/PreferencesTab'
import classes from './SettingsPage.module.css'

const items = [
  { key: 'categories', label: 'Categories', icon: 'bookmarks' },
  { key: 'integrations', label: 'Integrations', icon: 'extension' },
  { key: 'preferences', label: 'Preferences', icon: 'build' },
]

const Actions = () => (
  <a href="https://docs.readflow.app" rel="noreferrer noopener" target="_blank">
    Go to docs
  </a>
)

const PageHeader = () => (
  <>
    <Appbar>
      <Actions />
    </Appbar>
    <Header>
      <h1>Settings</h1>
      <AboutButton />
      <Tabs items={items} />
    </Header>
  </>
)

const SettingsPage = ({ match }: RouteComponentProps) => (
  <Page header={<PageHeader />} className={classes.settings}>
    <Switch>
      <Route path={match.path + '/categories'} component={CategoriesRoutes} />
      <Route path={match.path + '/integrations'} component={IntegrationsRoutes} />
      <Route exact path={match.path + '/preferences'} component={PreferencesTab} />
      <Redirect exact from={match.path} to={match.path + '/categories'} />
    </Switch>
  </Page>
)

export default SettingsPage
