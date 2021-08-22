import React, { FC } from 'react'

import { Panel } from '..'
import { Page } from '../../layout'
import { Button, ErrorPanel } from '..'
import { Link } from 'react-router-dom'

interface Props {
  title?: string
}

export const ErrorPage: FC<Props> = ({ title = 'Error', children }) => (
  <Page title={title}>
    <Panel>
      <ErrorPanel
        title={title}
        actions={
          <Button title="Back to home page" as={Link} to="/" variant="danger">
            Dismiss
          </Button>
        }
      >
        {children}
      </ErrorPanel>
    </Panel>
  </Page>
)
