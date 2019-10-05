import React from 'react'
import { Redirect, Route, RouteComponentProps, Switch } from 'react-router'

import ButtonIcon from '../components/ButtonIcon'
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
import NotificationButton from './components/NotificationButton'
import Tabs from './components/Tabs'
import PreferencesTab from './preferences/PreferencesTab'
import AddRuleForm from './rules/AddRuleForm'
import EditRuleTab from './rules/EditRuleTab'
import RulesTab from './rules/RulesTab'

const items = [
  { key: 'categories', label: 'Categories', icon: 'bookmarks' },
  { key: 'rules', label: 'Rules', icon: 'directions' },
  { key: 'api-keys', label: 'API keys', icon: 'verified_users' },
  { key: 'archive-services', label: 'Archive service', icon: 'backup' },
  { key: 'preferences', label: 'Preferences', icon: 'build' }
]

type AllProps = RouteComponentProps<{}>

const Actions = () => (
  <>
    <a href="https://about.readflow.app/docs/en/" rel="noreferrer noopener" target="_blank">
      Go to docs
    </a>
    <NotificationButton />
  </>
)

const PageHeader = () => (
  <>
    <Appbar actions={<Actions />} />
    <Header>
      <h1>Settings</h1>
      <ButtonIcon icon="info" to="/about" title="About" />
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
      <Route exact path={match.path + '/rules'} component={RulesTab} />
      <Route exact path={match.path + '/rules/add'} component={AddRuleForm} />
      <Route exact path={match.path + '/rules/:id'} component={EditRuleTab} />
      <Route exact path={match.path + '/preferences'} component={PreferencesTab} />
      <Redirect exact from={match.path} to={match.path + '/categories'} />
    </Switch>
  </Page>
)
