import React from 'react'
import ReactModal from "react-modal"

import ButtonIcon from '../../common/ButtonIcon'
import styles from '../../common/Dialog.module.css'

import useKeyboard from '../../hooks/useKeyboard'
import { useModal } from 'react-modal-hook'
import AddArticleForm from './AddArticleForm'
import { Category } from '../../categories/models'

type Props = {
  category?: Category
}

export default ({category}: Props) => {
  const [showAddModal, hideAddModal] = useModal(
    () => (
      <ReactModal
        isOpen
        shouldCloseOnEsc
        shouldCloseOnOverlayClick
        shouldFocusAfterRender
        appElement={document.getElementById('root')!}
        onRequestClose={hideAddModal}
        className={styles.dialog}
        overlayClassName={styles.overlay}>
        <AddArticleForm category={category} onSuccess={hideAddModal} onCancel={hideAddModal}/>
      </ReactModal>
    )
  )

  useKeyboard('+', showAddModal)

  return (
    <ButtonIcon
      title="Add new article"
      onClick={showAddModal}
      floating
      icon="add"
      primary />
  )
}
