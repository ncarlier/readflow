import React, { useContext, useState } from 'react'
import { Link, useRouteMatch } from 'react-router-dom'

import { MessageContext } from '../../../contexts/MessageContext'
import { ErrorPanel, Logo, Button } from '../../../components'
import DeleteOutgoingWebhooksButton from './DeleteOutgoingWebhookButton'
import OutgoingWebhooksList from './OutgoingWebhookList'

export default () => {
  const { showMessage } = useContext(MessageContext)
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [selection, setSelection] = useState<number[]>([])
  const { path } = useRouteMatch()

  const onDeleteSuccessHandler = (msg: string) => {
    showMessage(msg)
    setSelection([])
  }

  return (
    <section>
      <header>
        <h2>
          <Logo name="webhook" style={{ maxWidth: '2em', verticalAlign: 'middle' }} />
          Outgoing Webhooks
        </h2>
        <DeleteOutgoingWebhooksButton
          selection={selection}
          onError={setErrorMessage}
          onSuccess={onDeleteSuccessHandler}
        />
        <Button
          title="Add new outgoing webhook"
          as={Link}
          to={{
            pathname: path + 'outgoing-webhooks/add',
            state: { modal: true },
          }}
        >
          Add
        </Button>
      </header>
      <p>Outgoing webhooks allow external integration to receive articles.</p>
      {errorMessage != null && <ErrorPanel title="Unable to delete outgoing webhook(s)">{errorMessage}</ErrorPanel>}
      <OutgoingWebhooksList onSelected={setSelection} />
    </section>
  )
}
