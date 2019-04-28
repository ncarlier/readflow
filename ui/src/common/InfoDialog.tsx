/* eslint-disable @typescript-eslint/no-non-null-assertion */
import React, { ReactNode } from 'react'
import ReactModal from 'react-modal'

import Button from './Button'
import styles from './Dialog.module.css'
import Panel from './Panel'

interface Props {
  title: string
  children: ReactNode
  onOk: (e: any) => void
}

export default ({ title, children, onOk }: Props) => (
  <ReactModal
    isOpen
    shouldCloseOnEsc
    shouldCloseOnOverlayClick
    shouldFocusAfterRender
    appElement={document.getElementById('root')!}
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
        <Button primary onClick={onOk}>
          Ok
        </Button>
      </footer>
    </Panel>
  </ReactModal>
)
