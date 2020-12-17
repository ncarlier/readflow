import React, { useContext, useState } from 'react'
import { Link } from 'react-router-dom'

import Button from '../../../components/Button'
import { OnSelectedFn } from '../../../components/DataTable'
import { MessageContext } from '../../../context/MessageContext'
import ErrorPanel from '../../../error/ErrorPanel'
import DeleteOutboundServicesButton from './DeleteOutgoingWebhookButton'
import OutboundServicesList from './OutgoingWebhookList'

export default () => {
  const { showMessage } = useContext(MessageContext)
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [selection, setSelection] = useState<number[]>([])

  const onSelectedHandler: OnSelectedFn = (ids) => {
    setSelection(ids)
  }

  return (
    <section>
      <header>
        <h2>Outgoing Webhooks</h2>
        <DeleteOutboundServicesButton selection={selection} onError={setErrorMessage} onSuccess={showMessage} />
        <Button
          title="Add new outgoing webhook"
          as={Link}
          to={{
            pathname: '/outbound',
            state: { modal: true },
          }}
        >
          Add
        </Button>
      </header>
      <p>Outgoing webhooks allow external integration to receive articles.</p>
      {errorMessage != null && <ErrorPanel title="Unable to delete outgoing webhook(s)">{errorMessage}</ErrorPanel>}
      <OutboundServicesList onSelected={onSelectedHandler} />
    </section>
  )
}
