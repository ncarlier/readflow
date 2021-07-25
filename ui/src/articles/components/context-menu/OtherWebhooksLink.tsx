import React from 'react'
import ReactModal from 'react-modal'
import { useModal } from 'react-modal-hook'

import LinkIcon from '../../../components/LinkIcon'
import styles from '../../../components/Dialog.module.css'
import { OutgoingWebhook } from '../../../settings/intergrations/outgoing-webhook/models'
import OutgoingWebhooksPanel from './OutgoingWebhooksPanel'

interface Props {
  webhooks: OutgoingWebhook[]
  sendArticle: (alias: string) => any
}

const OtherWebhooksLink = ({ webhooks, sendArticle }: Props) => {
  const [showSenToModal, hideSendToModal] = useModal(() => (
    <ReactModal
      isOpen
      shouldCloseOnEsc
      shouldCloseOnOverlayClick
      shouldFocusAfterRender
      appElement={document.getElementById('root')!}
      onRequestClose={hideSendToModal}
      className={styles.dialog}
      overlayClassName={styles.overlay}
      style={{ content: { minWidth: '50vw' } }}
    >
      <OutgoingWebhooksPanel onCancel={hideSendToModal} sendArticle={sendArticle} webhooks={webhooks} />
    </ReactModal>
  ))

  if (webhooks.length) {
    return (
      <LinkIcon title="Send article to ..." icon="backup" onClick={showSenToModal}>
        <span>Send to ...</span>
      </LinkIcon>
    )
  }
  return null
}

export default OtherWebhooksLink
