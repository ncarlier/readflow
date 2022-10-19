/* eslint-disable @typescript-eslint/no-non-null-assertion */
import React, { FC, PropsWithChildren } from 'react'
import ReactModal from 'react-modal'

import { Button, Panel } from '.'
import styles from './Dialog.module.css'

interface Props extends PropsWithChildren {
  title: string
  confirmLabel?: string
  onConfirm: (e: any) => void
  onCancel?: (e: any) => void
}

export const ConfirmDialog: FC<Props> = ({ title, confirmLabel = 'ok', children, onConfirm, onCancel }) => (
  <ReactModal
    isOpen
    shouldCloseOnEsc
    shouldCloseOnOverlayClick
    shouldFocusAfterRender
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
