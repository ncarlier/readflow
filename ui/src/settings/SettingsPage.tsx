import React from 'react'
import { Switch, Route, RouteComponentProps, Redirect } from 'react-router'

import Page from '../common/Page'
import Header from './components/Header'
import Content from '../common/Content'
import Tabs from './components/Tabs'
import CategoriesTab from './categories/CategoriesTab'
import ApiKeysTab from './api-keys/ApiKeysTab'
import ArchiveServiceTab from './archive-services/ArchiveServicesTab'
import ButtonIcon from '../common/ButtonIcon'
import AddCategoryForm from './categories/AddCategoryForm'
import EditCategoryTab from './categories/EditCategoryTab'
import AddApiKeyForm from './api-keys/AddApiKeyForm'
import EditApiKeyTab from './api-keys/EditApiKeyTab'
import AddArchiveServiceForm from './archive-services/AddArchiveServiceForm'
import EditArchiveServiceTab from './archive-services/EditArchiveServiceTab'
import RulesTab from './rules/RulesTab'
import AddRuleForm from './rules/AddRuleForm'
import EditRuleTab from './rules/EditRuleTab'

const items = [
  { key: 'categories', label: 'Categories', icon: 'bookmarks'},
  { key: 'rules', label: 'Rules', icon: 'directions' },
  { key: 'api-keys', label: 'API keys', icon: 'verified_users' },
  { key: 'archive-services', label: 'Archive service', icon: 'backup' },
]

type AllProps = RouteComponentProps<{}>

export default ({match}: AllProps) => (
  <Page title="">
    <Header>
      <h1>Settings</h1>
      <ButtonIcon icon="info" to="/about" title="About" />
      <Tabs items={items} />
    </Header>
    <Content>
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
        <Redirect exact from={match.path} to={match.path + '/categories'} />
      </Switch>
    </Content>
  </Page>
)
