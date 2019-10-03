import React, { ReactNode } from 'react'

import Button from '../components/Button'
import Panel from '../components/Panel'
import Page from '../layout/Page'
import ErrorPanel from './ErrorPanel'

interface Props {
  title?: string
  children: ReactNode
}

export default ({ title = 'Error', children }: Props) => (
  <Page title={title}>
    <Panel>
      <ErrorPanel
        title={title}
        actions={
          <Button title="Back to home page" to="/" danger>
            Dismiss
          </Button>
        }
      >
        {children}
      </ErrorPanel>
    </Panel>
  </Page>
)
