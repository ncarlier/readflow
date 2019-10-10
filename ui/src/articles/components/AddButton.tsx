/* eslint-disable @typescript-eslint/no-non-null-assertion */
import React from 'react'
import ReactModal from 'react-modal'
import { useModal } from 'react-modal-hook'

import { Category } from '../../categories/models'
import ButtonIcon from '../../components/ButtonIcon'
import styles from '../../components/Dialog.module.css'
import useKeyboard from '../../hooks/useKeyboard'
import { Article } from '../models'
import AddArticleForm from './AddArticleForm'

interface Props {
  category?: Category
  onSuccess: (article: Article) => void
}

export default ({ category, onSuccess }: Props) => {
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
      <AddArticleForm
        category={category}
        onSuccess={a => {
          onSuccess(a)
          hideAddModal()
        }}
        onCancel={hideAddModal}
      />
    </ReactModal>
  ))

  useKeyboard('+', showAddModal)

  return <ButtonIcon title="Add new article" onClick={showAddModal} floating icon="add" variant="primary" />
}
