import React from 'react'

import { IncomingWebhook } from './models'
import Bookmarklet from './Bookmarklet'
import { HelpLink, Logo, CopyableField } from '../../../components'
import HelpSection from '../../HelpSection'
import QRCodeIncomingWebhookButton from './QRCodeIncomingWebhookButton'
import { getAPIURL } from '../../../helpers'

interface Props {
  data?: IncomingWebhook
}

export default ({ data }: Props) => (
  <HelpSection>
    <Logo name="webhook" style={{ maxWidth: '10%', verticalAlign: 'middle' }} />
    <span>
      Use incoming webhooks to post articles to your Readflow.<br />
      Messages are sent via an HTTP POST request to Readflow integration URL.<br />
      You can customize your integration with a script.
      {data && (
        <table>
          <tbody>
            <tr>
              <th>Ingestion URL</th>
              <td>{getAPIURL('/articles')}</td>
            </tr>
            {data.email &&
              <tr>
                <th>Ingestion email</th>
                <td>
                  <CopyableField value={data.email} />
                </td>
              </tr>
            }
            <tr>
              <th>Token</th>
              <td>
                <CopyableField masked value={data.token} />
              </td>
            </tr>
            <tr>
              <th>Bookmarklet</th>
              <td>
                <Bookmarklet token={data.token} />
              </td>
            </tr>
            <tr>
              <th>QR code</th>
              <td>
                <QRCodeIncomingWebhookButton token={data.token} />
              </td>
            </tr>
          </tbody>
        </table>
      )}
    </span>
    <HelpLink href="https://docs.readflow.app/en/integrations/incoming-webhook/">Help</HelpLink>
  </HelpSection>
)
