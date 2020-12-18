import React from 'react'

import { API_BASE_URL } from '../../../constants'
import Masked from '../../../components/Masked'
import { IncomingWebhook } from './models'
import Bookmarklet from './Bookmarklet'
import HelpLink from '../../../components/HelpLink'
import WebhookLogo from '../WebhookLogo'
import HelpSection from '../../HelpSection'

interface Props {
  data?: IncomingWebhook
}

export default ({ data }: Props) => (
  <HelpSection>
    <WebhookLogo maxWidth="10%" />
    <span>
      Use incoming webhooks to post articles to your Readflow. <br />
      Messages are sent via an HTTP POST request to Readflow integration URL.
      {data && (
        <table>
          <tbody>
            <tr>
              <th>Ingestion URL</th>
              <td>{API_BASE_URL + '/articles'}</td>
            </tr>
            <tr>
              <th>Token</th>
              <td>
                <Masked value={data.token} />
              </td>
            </tr>
            <tr>
              <th>Bookmarklet</th>
              <td>
                <Bookmarklet token={data.token} />
              </td>
            </tr>
          </tbody>
        </table>
      )}
    </span>
    <HelpLink href="https://about.readflow.app/docs/en/third-party/create/integration-api/">Help</HelpLink>
  </HelpSection>
)
