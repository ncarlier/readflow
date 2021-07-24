import React, { useContext, useState } from 'react'
import { Link, useRouteMatch } from 'react-router-dom'

import Button from '../../../components/Button'
import { MessageContext } from '../../../context/MessageContext'
import ErrorPanel from '../../../error/ErrorPanel'
import Logo from '../../../logos/Logo'
import DeleteIncomingWebhookButton from './DeleteIncomingWebhookButton'
import IncomingWebhookList from './IncomingWebhookList'

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
          Incoming Webhooks
        </h2>
        <DeleteIncomingWebhookButton
          selection={selection}
          onError={setErrorMessage}
          onSuccess={onDeleteSuccessHandler}
        />
        <Button
          id="add-new-incoming-webhook"
          title="Add new incoming webhook"
          as={Link}
          to={{
            pathname: path + 'incoming-webhooks/add',
            state: { modal: true },
          }}
        >
          Add
        </Button>
      </header>
      <p>Incoming webhooks allow external integrations to send articles.</p>
      {errorMessage != null && <ErrorPanel title="Unable to delete incoming webhook(s)">{errorMessage}</ErrorPanel>}
      <IncomingWebhookList onSelected={setSelection} />
    </section>
  )
}
