/* eslint-disable @typescript-eslint/no-non-null-assertion */
import React, { ReactNode } from 'react'
import ReactModal from 'react-modal'

import Button from './Button'
import styles from './Dialog.module.css'
import Panel from './Panel'

interface Props {
  title: string
  confirmLabel: string
  children: ReactNode
  onConfirm: (e: any) => void
  onCancel?: (e: any) => void
}

export default ({ title, confirmLabel, children, onConfirm, onCancel }: Props) => (
  <ReactModal
    isOpen
    shouldCloseOnEsc
    shouldCloseOnOverlayClick
    shouldFocusAfterRender
    appElement={document.getElementById('root')!}
    onRequestClose={onCancel}
    className={styles.dialog}
    overlayClassName={styles.overlay}
  >
    <Panel>
      <header>
        <h1>{title}</h1>
      </header>
      <section>{children}</section>
      <footer>
        {onCancel && <Button onClick={onCancel}>Cancel</Button>}
        <Button id="modal-confirm" variant="primary" onClick={onConfirm} autoFocus>
          {confirmLabel}
        </Button>
      </footer>
    </Panel>
  </ReactModal>
)
