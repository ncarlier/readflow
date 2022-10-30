import React from 'react'
import ReactModal from 'react-modal'
import { useModal } from 'react-modal-hook'

import { EditArticleForm } from './components/EditArticleForm'
import { Article } from './models'
import styles from '../components/Dialog.module.css'

export const useArticleEditModal = (article: Article) => {
  const [showEditModal, hideEditModal] = useModal(
    () => (
      <ReactModal
        isOpen
        shouldCloseOnEsc
        shouldCloseOnOverlayClick
        shouldFocusAfterRender
        onRequestClose={hideEditModal}
        className={styles.dialog}
        overlayClassName={styles.overlay}
        style={{ content: { minWidth: '50vw' } }}
      >
        <EditArticleForm
          article={article}
          onSuccess={() => {
            hideEditModal()
          }}
          onCancel={hideEditModal}
        />
      </ReactModal>
    ),
    [article],
  )
  return [showEditModal, hideEditModal]
}
