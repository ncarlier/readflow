import React from 'react'
import { Redirect, Route, RouteComponentProps, Switch } from 'react-router'

import Appbar from '../layout/Appbar'
import Page from '../layout/Page'
import AboutButton from './about/AboutButton'
import CategoriesRoutes from './categories/routes'
import Header from './components/Header'
import Tabs from './components/Tabs'
import IntegrationsRoutes from './intergrations/routes'
import PreferencesTab from './preferences/PreferencesTab'
import classes from './SettingsPage.module.css'

const items = [
  { key: 'categories', label: 'Categories', icon: 'bookmarks' },
  { key: 'integrations', label: 'Integrations', icon: 'extension' },
  { key: 'preferences', label: 'Preferences', icon: 'build' },
]

const Actions = () => (
  <a href="https://about.readflow.app/docs/en/" rel="noreferrer noopener" target="_blank">
    Go to docs
  </a>
)

const PageHeader = () => (
  <>
    <Appbar actions={<Actions />} />
    <Header>
      <h1>Settings</h1>
      <AboutButton />
      <Tabs items={items} />
    </Header>
  </>
)

export default ({ match }: RouteComponentProps) => (
  <Page header={<PageHeader />} className={classes.settings}>
    <Switch>
      <Route path={match.path + '/categories'} component={CategoriesRoutes} />
      <Route path={match.path + '/integrations'} component={IntegrationsRoutes} />
      <Route exact path={match.path + '/preferences'} component={PreferencesTab} />
      <Redirect exact from={match.path} to={match.path + '/categories'} />
    </Switch>
  </Page>
)
