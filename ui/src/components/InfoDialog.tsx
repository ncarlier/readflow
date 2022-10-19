/* eslint-disable @typescript-eslint/no-non-null-assertion */
import React, { FC, PropsWithChildren } from 'react'
import ReactModal from 'react-modal'

import { Button, Panel } from '.'
import styles from './Dialog.module.css'

interface Props extends PropsWithChildren {
  title: string
  onOk: (e: any) => void
}

export const InfoDialog: FC<Props> = ({ title, children, onOk }) => (
  <ReactModal
    isOpen
    shouldCloseOnEsc
    shouldCloseOnOverlayClick
    shouldFocusAfterRender
    onRequestClose={onOk}
    className={styles.dialog}
    overlayClassName={styles.overlay}
  >
    <Panel>
      <header>
        <h1>{title}</h1>
      </header>
      <section>{children}</section>
      <footer>
        <Button variant="primary" onClick={onOk} autoFocus>
          Ok
        </Button>
      </footer>
    </Panel>
  </ReactModal>
)
