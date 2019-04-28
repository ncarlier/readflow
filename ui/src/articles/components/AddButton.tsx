/* eslint-disable @typescript-eslint/no-non-null-assertion */
import React from 'react'
import ReactModal from 'react-modal'
import { useModal } from 'react-modal-hook'

import { Category } from '../../categories/models'
import ButtonIcon from '../../common/ButtonIcon'
import styles from '../../common/Dialog.module.css'
import useKeyboard from '../../hooks/useKeyboard'
import AddArticleForm from './AddArticleForm'

interface Props {
  category?: Category
}

export default ({ category }: Props) => {
  const [showAddModal, hideAddModal] = useModal(() => (
    <ReactModal
      isOpen
      shouldCloseOnEsc
      shouldCloseOnOverlayClick
      shouldFocusAfterRender
      appElement={document.getElementById('root')!}
      onRequestClose={hideAddModal}
      className={styles.dialog}
      overlayClassName={styles.overlay}
    >
      <AddArticleForm category={category} onSuccess={hideAddModal} onCancel={hideAddModal} />
    </ReactModal>
  ))

  useKeyboard('+', showAddModal)

  return <ButtonIcon title="Add new article" onClick={showAddModal} floating icon="add" primary />
}
