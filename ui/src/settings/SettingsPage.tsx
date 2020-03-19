import React from 'react'
import { Redirect, Route, RouteComponentProps, Switch } from 'react-router'

import Appbar from '../layout/Appbar'
import Page from '../layout/Page'
import AddApiKeyForm from './api-keys/AddApiKeyForm'
import ApiKeysTab from './api-keys/ApiKeysTab'
import EditApiKeyTab from './api-keys/EditApiKeyTab'
import AddArchiveServiceForm from './archive-services/AddArchiveServiceForm'
import ArchiveServiceTab from './archive-services/ArchiveServicesTab'
import EditArchiveServiceTab from './archive-services/EditArchiveServiceTab'
import AddCategoryForm from './categories/AddCategoryForm'
import CategoriesTab from './categories/CategoriesTab'
import EditCategoryTab from './categories/EditCategoryTab'
import Header from './components/Header'
import Tabs from './components/Tabs'
import PreferencesTab from './preferences/PreferencesTab'
import AboutButton from './about/AboutButton'

const items = [
  { key: 'categories', label: 'Categories', icon: 'bookmarks' },
  { key: 'api-keys', label: 'API keys', icon: 'verified_users' },
  { key: 'archive-services', label: 'Archive service', icon: 'backup' },
  { key: 'preferences', label: 'Preferences', icon: 'build' }
]

type AllProps = RouteComponentProps<{}>

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

export default ({ match }: AllProps) => (
  <Page header={<PageHeader />}>
    <Switch>
      <Route exact path={match.path + '/categories'} component={CategoriesTab} />
      <Route exact path={match.path + '/categories/add'} component={AddCategoryForm} />
      <Route exact path={match.path + '/categories/:id'} component={EditCategoryTab} />
      <Route exact path={match.path + '/api-keys'} component={ApiKeysTab} />
      <Route exact path={match.path + '/api-keys/add'} component={AddApiKeyForm} />
      <Route exact path={match.path + '/api-keys/:id'} component={EditApiKeyTab} />
      <Route exact path={match.path + '/archive-services'} component={ArchiveServiceTab} />
      <Route exact path={match.path + '/archive-services/add'} component={AddArchiveServiceForm} />
      <Route exact path={match.path + '/archive-services/:id'} component={EditArchiveServiceTab} />
      <Route exact path={match.path + '/preferences'} component={PreferencesTab} />
      <Redirect exact from={match.path} to={match.path + '/categories'} />
    </Switch>
  </Page>
)
