import React from 'react'
import { Route, Switch, useRouteMatch } from 'react-router'
import AddCategoryForm from './AddCategoryForm'
import CategoriesTab from './CategoriesTab'
import EditCategoryTab from './EditCategoryTab'

export default () => {
  const { path } = useRouteMatch()
  return (
    <Switch>
      <Route exact path={path + '/'} component={CategoriesTab} />
      <Route exact path={path + '/add'} component={AddCategoryForm} />
      <Route exact path={path + '/:id'} component={EditCategoryTab} />
    </Switch>
  )
}
