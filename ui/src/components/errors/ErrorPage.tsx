import React, { FC, PropsWithChildren } from 'react'

import { Page } from '../../layout'
import { Button, ErrorPanel, Panel } from '..'
import { Link } from 'react-router-dom'

interface Props extends PropsWithChildren {
  title?: string
}

export const ErrorPage: FC<Props> = ({ title = 'Error', children }) => (
  <Page title={title}>
    <Panel>
      <ErrorPanel
        title={title}
        actions={
          <Button title="Back to home page" as={Link} replace to="/" variant="danger">
            Back
          </Button>
        }
      >
        {children}
      </ErrorPanel>
    </Panel>
  </Page>
)
