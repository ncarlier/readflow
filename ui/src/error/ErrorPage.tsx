import React, { ReactNode } from 'react'

import Page from '../common/Page'
import ErrorPanel from './ErrorPanel'
import Button from '../common/Button';
import Panel from '../common/Panel';

type Props = {
  title?: string
  children: ReactNode
}

export default ({title="Error", children}: Props) => (
  <Page title={title}>
    <Panel>
      <ErrorPanel title={title} actions={
        <Button title="Back to home page" to='/' danger>Dismiss</Button>
        }>
        {children}
      </ErrorPanel>
    </Panel>
  </Page>
)
