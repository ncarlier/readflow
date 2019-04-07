import React, { ReactNode } from "react"
import ReactModal from "react-modal"

import styles from './ConfirmDialog.module.css'
import Panel from "./Panel";
import Button from "./Button";
const customStyles = {
  content : {
    top         : '50%',
    left        : '50%',
    right       : 'auto',
    bottom      : 'auto',
    marginRight : '-50%',
    transform   : 'translate(-50%, -50%)'
  }
}

type Props = {
  title: string
  confirmLabel: string
  children: ReactNode
  onConfirm: (e: any) => void
  onCancel?: (e: any) => void
}

const ConfirmDialog = ({
  title,
  confirmLabel,
  children,
  onConfirm,
  onCancel,
}: Props) => (
  <ReactModal
    isOpen
    appElement={document.getElementById('root')!}
    className={styles.dialog}
    overlayClassName={styles.overlay}>
    <Panel>
      <header>
        <h1>{title}</h1>
      </header>
      <section>{children}</section>
      <footer>
        { onCancel && <Button onClick={onCancel}>Cancel</Button> }
        <Button primary onClick={onConfirm}>{confirmLabel}</Button>
      </footer>
    </Panel>
  </ReactModal>
)

export default ConfirmDialog;
